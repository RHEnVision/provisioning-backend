package background

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/rs/zerolog"
)

func dbCleanup(ctx context.Context, sleep time.Duration) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msgf("Started reservation cleanup %s", sleep.String())
	defer func() {
		logger.Debug().Msgf("Database reservation cleanup routine exited")
	}()

	ticker := time.NewTicker(sleep)

	cleanupReservations(ctx)

	for {
		select {
		case <-ticker.C:
			cleanupReservations(ctx)

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func cleanupReservations(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	sdao := dao.GetReservationDao(ctx)
	err := sdao.Cleanup(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error while performing reservation cleanup")
	}
}
