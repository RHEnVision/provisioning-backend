package jobs

import (
	"context"
	"time"

	"github.com/lzap/dejq"
)

type NoopJobArgs struct {
	AccountID     int64 `json:"account_id"`
	ReservationID int64 `json:"reservation_id"`
}

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
	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "No operation started")
	defer updateStatusAfter(ctx, args.ReservationID, "No operation finished", 1)

	// do nothing and return no error
	time.Sleep(5 * time.Second)

	return nil
}
