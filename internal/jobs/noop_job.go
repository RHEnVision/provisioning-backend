package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
)

type NoopJobArgs struct {
	AccountID     int64 `json:"account_id"`
	ReservationID int64 `json:"reservation_id"`

	// Returns an error when set, used only in tests.
	Fail bool `json:"fail"`
}

var NoopJobTask = queue.RegisterTask("noop", func(ctx context.Context, args NoopJobArgs) error {
	ctx = prepareContext(ctx, "noop", args, args.AccountID, args.ReservationID)
	err := HandleNoop(ctx, args)
	finishStep(ctx, args.ReservationID, err)
	return err
})

func EnqueueNoop(ctx context.Context, args NoopJobArgs) error {
	ctxLogger := ctxval.Logger(ctx)

	t := NoopJobTask.WithArgs(ctx, args)
	ctxLogger.Debug().Str("tid", t.ID).Msg("Adding noop job task")
	err := queue.JobQueue.Add(t)
	if err != nil {
		return fmt.Errorf("unable to enqueue task: %w", err)
	}

	return nil
}

var NoOperationFailure = errors.New("job failed on request")

func HandleNoop(ctx context.Context, args NoopJobArgs) error {
	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "No operation started")
	defer updateStatusAfter(ctx, args.ReservationID, "No operation finished", 1)

	// do nothing
	time.Sleep(10 * time.Millisecond)
	logger.Info().Msg("No operation finished")

	if args.Fail {
		return NoOperationFailure
	}

	return nil
}
