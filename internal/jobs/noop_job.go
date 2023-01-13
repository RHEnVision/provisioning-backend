package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

type NoopJobArgs struct {
	AccountID     int64
	ReservationID int64

	// Indicates that the test should fail, used only in tests.
	Fail bool
}

var NoOperationFailure = errors.New("job failed on request")

// Unmarshall arguments and handle error
func HandleNoop(ctx context.Context, job *worker.Job) {
	args, ok := job.Args.(NoopJobArgs)
	if !ok {
		ctxval.Logger(ctx).Error().Msgf("Type assertion error for job %s, unable to finish reservation: %#v", job.ID, job.Args)
		return
	}

	ctx = contextLogger(ctx, job, args.AccountID, args.ReservationID)

	jobErr := DoNoop(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
}

// Job logic, when error is returned the job status is updated accordingly
func DoNoop(ctx context.Context, args *NoopJobArgs) error {
	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "No operation started")
	defer updateStatusAfter(ctx, args.ReservationID, "No operation finished", 1)

	// skip job if reservation already contains errors
	err := checkExistingError(ctx, args.ReservationID)
	if err != nil {
		return fmt.Errorf("step skipped: %w", err)
	}

	// do nothing
	time.Sleep(10 * time.Millisecond)
	logger.Info().Msg("No operation finished")

	if args.Fail {
		return NoOperationFailure
	}

	return nil
}
