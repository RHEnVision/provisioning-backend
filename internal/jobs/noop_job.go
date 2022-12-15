package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/lzap/dejq"
)

type NoopJobArgs struct {
	AccountID     int64 `json:"account_id"`
	ReservationID int64 `json:"reservation_id"`

	// Indicates that the test should fail, used only in tests.
	Fail bool `json:"fail"`
}

var NoOperationFailure = errors.New("job failed on request")

// Unmarshall arguments and handle error
func HandleNoop(ctx context.Context, job dejq.Job) error {
	args := NoopJobArgs{}
	err := decodeJob(ctx, job, &args)
	if err != nil {
		return err
	}
	ctx = contextLogger(ctx, job.Type(), args, args.AccountID, args.ReservationID)

	jobErr := handleNoop(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
	return jobErr
}

// Job logic, when error is returned the job status is updated accordingly
func handleNoop(ctx context.Context, args *NoopJobArgs) error {
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
