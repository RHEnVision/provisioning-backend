package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
)

type NoopTwoArgs struct {
	AccountID     int64 `json:"account_id"`
	ReservationID int64 `json:"reservation_id"`

	// The first function (step) returns an error when set, used only in tests.
	Fail1 bool `json:"fail1"`
	// The second function (step) returns an error when set, used only in tests.
	Fail2 bool `json:"fail2"`
}

var NoopTwoTask = queue.RegisterTask("noop_two", func(ctx context.Context, args NoopTwoArgs) error {
	ctx = prepareContext(ctx, "noop_two", args, args.AccountID, args.ReservationID)

	err := HandleNoop1(ctx, args)
	finishStep(ctx, args.ReservationID, err)
	if err != nil {
		return err
	}

	err = HandleNoop2(ctx, args)
	finishStep(ctx, args.ReservationID, err)
	if err != nil {
		return err
	}

	return err
})

func EnqueueNoopTwo(ctx context.Context, args NoopTwoArgs) error {
	ctxLogger := ctxval.Logger(ctx)

	t := NoopTwoTask.WithArgs(ctx, args)
	ctxLogger.Debug().Str("tid", t.ID).Msg("Adding noop job task")
	err := queue.JobQueue.Add(t)
	if err != nil {
		return fmt.Errorf("unable to enqueue task: %w", err)
	}

	return nil
}

func HandleNoop1(ctx context.Context, args NoopTwoArgs) error {
	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Operation one started")
	defer updateStatusAfter(ctx, args.ReservationID, "Operation one finished", 1)

	// do nothing
	time.Sleep(10 * time.Millisecond)
	logger.Info().Msg("No operation finished")

	if args.Fail1 {
		return NoOperationFailure
	}

	return nil
}

func HandleNoop2(ctx context.Context, args NoopTwoArgs) error {
	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Operation two started")
	defer updateStatusAfter(ctx, args.ReservationID, "Operation two finished", 1)

	// do nothing
	time.Sleep(10 * time.Millisecond)
	logger.Info().Msg("No operation finished")

	if args.Fail2 {
		return NoOperationFailure
	}

	return nil
}
