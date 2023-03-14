package jobs

import (
	"context"
	"errors"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

type NoopJobArgs struct {
	ReservationID int64
	Fail          bool          // Fail forcefully (used in tests)
	Sleep         time.Duration // Sleep (delay) duration (used in tests)
}

var NoOperationFailure = errors.New("job failed on request")

// Unmarshall arguments and handle error
func HandleNoop(ctx context.Context, job *worker.Job) {
	args, ok := job.Args.(NoopJobArgs)
	if !ok {
		ctxval.Logger(ctx).Error().Msgf("Type assertion error for job %s, unable to finish reservation: %#v", job.ID, job.Args)
		return
	}

	logger := ctxval.Logger(ctx).With().Int64("reservation_id", args.ReservationID).Logger()
	ctx = ctxval.WithLogger(ctx, &logger)

	jobErr := DoNoop(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
}

// Job logic, when error is returned the job status is updated accordingly
func DoNoop(ctx context.Context, args *NoopJobArgs) error {
	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "No operation started")
	defer updateStatusAfter(ctx, args.ReservationID, "No operation finished", 1)

	// do nothing
	_ = sleepCtx(ctx, args.Sleep)
	logger.Info().Msg("No operation finished")

	if args.Fail {
		return NoOperationFailure
	}

	return nilUnlessTimeout(ctx)
}
