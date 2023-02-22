package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/getsentry/sentry-go"
)

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
	logger := ctxval.Logger(ctx)

	rDao := dao.GetReservationDao(ctx)

	reservation, err := rDao.GetById(ctx, reservationId)
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update job status: get by id")
		return
	}
	logger.Debug().Msgf("Job step: %d/%d", reservation.Step, reservation.Steps)

	// if this was the last step, set the success flag
	if reservation.Step >= reservation.Steps {
		logger.Info().Msgf("All jobs executed, marking job as success")
		err = rDao.FinishWithSuccess(ctx, reservationId)
		if err != nil {
			logger.Warn().Err(err).Msg("unable to update job status: finish")
		}
	}
}

func finishWithError(ctx context.Context, reservationId int64, jobError error) {
	logger := ctxval.Logger(ctx)
	logger.Warn().Err(jobError).Msgf("Job returned an error: %s", jobError.Error())

	// Send the error to Sentry
	if config.Sentry.Enabled {
		sentry.CaptureException(jobError)
	}

	rDao := dao.GetReservationDao(ctx)
	err := rDao.FinishWithError(ctx, reservationId, jobError.Error())
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update job status: finish")
	}
}

func updateStatusBefore(ctx context.Context, id int64, status string) {
	updateStatusAfter(ctx, id, status, 0)
}

func updateStatusAfter(ctx context.Context, id int64, status string, addSteps int) {
	logger := ctxval.Logger(ctx)
	logger.Debug().Bool("step", true).Msgf("Reservation status change: '%s'", status)
	if addSteps != 0 {
		logger.Trace().Bool("step", true).Msgf("Increased step number by: %d", addSteps)
	}

	rDao := dao.GetReservationDao(ctx)

	err := rDao.UpdateStatus(ctx, id, status, int32(addSteps))
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update step number: update")
	}
}

// nolint: goerr113
func checkExistingError(ctx context.Context, reservationId int64) error {
	resDao := dao.GetReservationDao(ctx)
	reservation, err := resDao.GetById(ctx, reservationId)
	if err != nil {
		return fmt.Errorf("cannot find reservation: %w", err)
	}
	if reservation.Error != "" {
		ctxval.Logger(ctx).Warn().Msg("Reservation already contains error, skipping job")
		return errors.New(reservation.Error)
	}

	return nil
}
