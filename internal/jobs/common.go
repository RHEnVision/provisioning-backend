package jobs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
)

// prepareContext puts account ID into the context and enriches context logger with reservation ID.
func prepareContext(ctx context.Context, jobName string, args interface{}, accountId, reservationId int64) context.Context {
	newContext := ctxval.WithAccountId(ctx, accountId)
	logger := ctxval.Logger(newContext).With().
		Int64("reservation_id", reservationId).
		Int64("account_id", accountId).Logger()
	logger.Info().Interface("args", args).Msgf("Processing job: '%s'", jobName)
	newContext = ctxval.WithLogger(newContext, &logger)
	return newContext
}

// finishStep sets reservation success and error field accordingly.
func finishStep(ctx context.Context, reservationId int64, jobErr error) {
	if jobErr != nil {
		// TODO: increase counter of failed jobs in Prometheus
		finishWithError(ctx, reservationId, jobErr)
	} else {
		// TODO: increase counter of successful jobs in Prometheus
		finishWithSuccess(ctx, reservationId)
	}
}

func finishWithSuccess(ctx context.Context, reservationId int64) {
	ctxLogger := ctxval.Logger(ctx)
	rDao := dao.GetReservationDao(ctx)

	reservation, err := rDao.GetById(ctx, reservationId)
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update job status: get by id")
		return
	}
	ctxLogger.Debug().Msgf("Job step: %d/%d", reservation.Step, reservation.Steps)

	// if this was the last step, set the success flag
	if reservation.Step >= reservation.Steps {
		ctxLogger.Info().Msgf("All jobs executed, marking job as success")
		err = rDao.FinishWithSuccess(ctx, reservationId)
		if err != nil {
			ctxLogger.Warn().Err(err).Msg("unable to update job status: finish")
		}
	}
}

func finishWithError(ctx context.Context, reservationId int64, jobError error) {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Warn().Err(jobError).Msgf("Job returned an error: %s", jobError.Error())

	rDao := dao.GetReservationDao(ctx)
	err := rDao.FinishWithError(ctx, reservationId, jobError.Error())
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update job status: finish")
	}
}

// updateStatusBefore sets the status string.
func updateStatusBefore(ctx context.Context, id int64, status string) {
	updateStatusAfter(ctx, id, status, 0)
}

// updateStatusAfter sets the status string and modifies step counter.
func updateStatusAfter(ctx context.Context, id int64, status string, addSteps int) {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Bool("step", true).Msgf("Reservation status change: '%s'", status)
	if addSteps != 0 {
		ctxLogger.Trace().Bool("step", true).Msgf("Increased step number by: %d", addSteps)
	}

	rDao := dao.GetReservationDao(ctx)

	err := rDao.UpdateStatus(ctx, id, status, int32(addSteps))
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update step number: update")
	}
}
