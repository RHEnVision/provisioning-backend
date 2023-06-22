package background

import (
	"context"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/rs/zerolog"
)

func dbStatsLoop(ctx context.Context, sleep time.Duration) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msgf("Started database statistics routine")
	defer func() {
		logger.Debug().Msgf("Database statistics routine exited")
	}()
	ticker := time.NewTicker(sleep)

	for {
		select {
		case <-ticker.C:
			metrics.ObserveDbStatsDuration(func() {
				err := dbStatsTick(ctx)
				if err != nil {
					logger.Error().Err(err).Msg("Error while performing database statistics query")
				}
			})

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func dbStatsTick(ctx context.Context) error {
	sdao := dao.GetStatDao(ctx)
	stats, err := sdao.Get(ctx)
	if err != nil {
		return fmt.Errorf("stats error: %w", err)
	}

	for _, s := range stats.Usage24h {
		metrics.SetReservations24hCount(s.Result, s.Provider, s.Count)
	}
	for _, s := range stats.Usage28d {
		metrics.SetReservations28dCount(s.Result, s.Provider, s.Count)
	}

	return nil
}
