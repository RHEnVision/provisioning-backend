package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/notifications"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/rs/zerolog"
)

var (
	ErrTypeAssertion = errors.New("type assert error")
	ErrPanicInJob    = errors.New("panic during job")
)

func finishJob(ctx context.Context, reservationId int64, jobErr error) {
	nc := notifications.GetNotificationClient(ctx)

	if jobErr != nil {
		nc.FailedLaunch(ctx, reservationId, jobErr)
		finishWithError(ctx, reservationId, jobErr)
	} else {
		nc.SuccessfulLaunch(ctx, reservationId)
		finishWithSuccess(ctx, reservationId)
	}
}

func finishWithSuccess(ctx context.Context, reservationId int64) {
	logger := zerolog.Ctx(ctx)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		// the original context is expired and unusable at this point
		ctx = copyContext(ctx)
	}

	// get details
	rDao := dao.GetReservationDao(ctx)
	reservation, err := rDao.GetById(ctx, reservationId)
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update job status: get by id")
		return
	}
	if reservation.Step == reservation.Steps {
		logger.Info().Msgf("Finishing reservation with success at step %d/%d", reservation.Step, reservation.Steps)
	} else {
		logger.Error().Msgf("Finishing reservation with success at step %d/%d", reservation.Step, reservation.Steps)
	}

	// total count of reservations
	metrics.IncReservationCount(reservation.Provider.String(), "success")

	// and finish
	err = rDao.FinishWithSuccess(ctx, reservationId)
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update job status: finish")
	}
}

// finishWithError closes a reservation and sets it into error state. Error message is also
// stored into the reservation.
func finishWithError(ctx context.Context, reservationId int64, jobError error) {
	logger := zerolog.Ctx(ctx)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		// the original context is expired and unusable at this point
		ctx = copyContext(ctx)
	}

	// get details
	rDao := dao.GetReservationDao(ctx)
	reservation, err := rDao.GetById(ctx, reservationId)
	if err != nil {
		logger.Warn().Err(err).Msg("unable to update job status: get by id")
		return
	}
	logger.Error().Err(jobError).Msgf("Finishing reservation with error at step %d/%d", reservation.Step, reservation.Steps)

	// total count of reservations
	metrics.IncReservationCount(reservation.Provider.String(), "failure")

	// and finish
	err = rDao.FinishWithError(ctx, reservationId, jobError.Error())
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
	logger := zerolog.Ctx(ctx)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		status = "Timeout"
		// the original context is expired and unusable at this point
		ctx = copyContext(ctx)
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
//
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

var ErrTryAgain = errors.New("try again")

// wait and keep retrying until function returns nil in given millisecond intervals
//
//nolint:wrapcheck
func waitAndRetry(ctx context.Context, f func() error, intervalsMs ...int) error {
	if len(intervalsMs) < 2 {
		panic("number of retries must be 2 or more")
	}

	var err error
	retries := 0
	for _, interval := range intervalsMs {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		retries += 1
		err = f()
		if err == nil || ctx.Err() != nil {
			break
		}
	}

	logger := zerolog.Ctx(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Error when retrying")
	}
	logger.Trace().Int("retries", retries-1).Msgf("Number of tries %d out of %d", retries, len(intervalsMs))

	return err
}
