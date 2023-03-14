package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
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

// finishWithError closes a reservation and sets it into error state. Error message is also
// stored into the reservation.
func finishWithError(ctx context.Context, reservationId int64, jobError error) {
	logger := ctxval.Logger(ctx)
	logger.Warn().Err(jobError).Msgf("Job returned an error: %s", jobError.Error())
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		ctx = context.Background()
		ctxval.Logger(ctx).Warn().Msgf("Updating reservation %d after context deadline", reservationId)
	}

	rDao := dao.GetReservationDao(ctx)
	err := rDao.FinishWithError(ctx, reservationId, jobError.Error())
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update job status: finish")
	}
}

// updateStatusBefore is called after every step function within a job. It updates reservation status
// message.
func updateStatusBefore(ctx context.Context, id int64, status string) {
	updateStatusAfter(ctx, id, status, 0)
}

// updateStatusAfter is called after every step function within a job. It updates reservation status
// message and step counter. When context deadline was exceeded, it sets the status message to "Timeout".
func updateStatusAfter(ctx context.Context, id int64, status string, addSteps int) {
	logger := ctxval.Logger(ctx)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		status = "Timeout"
		ctx = context.Background()
		ctxval.Logger(ctx).Warn().Msgf("Updating reservation %d after context deadline", id)
	}

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

func nilUnlessTimeout(ctx context.Context) error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("context timeout: %w", ctx.Err())
	}
	return nil
}

// sleep with context deadline
//nolint:wrapcheck
func sleepCtx(ctx context.Context, d time.Duration) error {
	afterCh := time.After(d)
	select {
	case <-afterCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
