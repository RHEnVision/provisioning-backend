package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/lzap/dejq"
)

type NoopJobArgs struct {
	AccountID     int64 `json:"account_id"`
	ReservationID int64 `json:"reservation_id"`
}

func EnqueueNoop(ctx context.Context, args *NoopJobArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Interface("arg", args).Msgf("Enqueuing no operation job: %+v", args)

	pj := dejq.PendingJob{
		Type: TypeNoop,
		Body: args,
	}
	err := Queue.Enqueue(ctx, pj)
	if err != nil {
		return fmt.Errorf("unable to enqueue: %w", err)
	}

	return nil
}

func HandleNoop(ctx context.Context, job dejq.Job) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started no operation job")

	args := NoopJobArgs{}
	err := job.Decode(&args)
	if err != nil {
		ctxLogger.Error().Err(err).Msg("unable to decode arguments")
		return fmt.Errorf("unable to decode args: %w", err)
	}
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing no operation job")

	// do nothing and update status
	rDao, err := dao.GetReservationDao(ctx)
	if err != nil {
		return fmt.Errorf("cannot get reservation DAO: %w", err)
	}
	err = rDao.Finish(ctx, args.ReservationID, true, "Finished")
	if err != nil {
		return fmt.Errorf("cannot update reservation status: %w", err)
	}

	return nil
}
