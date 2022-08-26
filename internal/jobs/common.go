package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/lzap/dejq"
)

func decodeJob(ctx context.Context, job dejq.Job, argInterface interface{}) error {
	err := job.Decode(&argInterface)
	if err != nil {
		ctxval.Logger(ctx).Error().Err(err).Msgf("Unable to decode arguments for job '%s'", job.Type())
		return fmt.Errorf("decode error: %w", err)
		// TODO: increase counter of failed jobs in Prometheus
	}
	return nil
}

func contextLogger(ctx context.Context, jobName string, args interface{}, accountId, reservationId int64) context.Context {
	newContext := ctxval.WithAccountId(ctx, accountId)
	logger := ctxval.Logger(newContext).With().Int64("reservation_id", reservationId).Logger()
	logger.Info().Interface("args", args).Msgf("Processing job: '%s'", jobName)
	newContext = ctxval.WithLogger(newContext, &logger)
	return newContext
}

func finishJob(ctx context.Context, reservationId int64, jobErr error) {
	if jobErr != nil {
		// TODO: increase counter of failed jobs in Prometheus
		finishWithError(ctx, reservationId, jobErr)
	} else {
		// TODO: increase counter of successful jobs in Prometheus
		finishWithSuccess(ctx, reservationId)
	}
}

func finishWithSuccess(ctx context.Context, reservationId int64) {
	ctxLogger := ctxval.Logger(ctx).With().Int64("reservation_id", reservationId).Logger()

	rDao, err := dao.GetReservationDao(ctx)
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update job status: dao")
		return
	}

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
	ctxLogger := ctxval.Logger(ctx).With().Int64("reservation_id", reservationId).Logger()
	ctxLogger.Warn().Err(jobError).Msgf("Job returned an error: %s", jobError.Error())

	rDao, err := dao.GetReservationDao(ctx)
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update job status: dao")
		return
	}
	err = rDao.FinishWithError(ctx, reservationId, jobError.Error())
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update job status: finish")
	}
}

func updateStatusBefore(ctx context.Context, id int64, status string) {
	updateStatusAfter(ctx, id, status, 0)
}

func updateStatusAfter(ctx context.Context, id int64, status string, addSteps int) {
	ctxLogger := ctxval.Logger(ctx).With().Int64("reservation_id", id).Logger()
	ctxLogger.Debug().Bool("step", true).Msgf("Reservation status change: '%s'", status)
	if addSteps != 0 {
		ctxLogger.Trace().Bool("step", true).Msgf("Increased step number by: %d", addSteps)
	}

	rDao, err := dao.GetReservationDao(ctx)
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update step number: dao")
		return
	}

	err = rDao.UpdateStatus(ctx, id, status, int32(addSteps))
	if err != nil {
		ctxLogger.Warn().Err(err).Msg("unable to update step number: update")
	}
}
