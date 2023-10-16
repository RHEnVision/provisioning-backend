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
	logger.Debug().Msgf("Started database statistics routine with tick interval %.2f seconds", sleep.Seconds())
	defer func() {
		logger.Debug().Msgf("Database statistics routine exited")
	}()
	ticker := time.NewTicker(sleep)

	// run one tick immediately to prevent prometheus gaps
	dbStatsObserveTick(ctx)

	for {
		select {
		case <-ticker.C:
			dbStatsObserveTick(ctx)

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func dbStatsObserveTick(ctx context.Context) {
	logger := zerolog.Ctx(ctx)
	metrics.ObserveDbStatsDuration(func() {
		err := dbStatsTick(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("Error while performing database statistics query")
		}
	})
}

func dbStatsTick(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	sdao := dao.GetStatDao(ctx)
	stats, err := sdao.Get(ctx, 10)
	if err != nil {
		return fmt.Errorf("stats error: %w", err)
	}

	var success, failure, pending int64
	for _, s := range stats.Usage24h {
		if s.Result == "success" {
			success += s.Count
		}
		if s.Result == "failure" {
			failure += s.Count
		}
		if s.Result == "pending" {
			pending += s.Count
		}
		metrics.SetReservations24hCount(s.Result, s.Provider, s.Count)
	}
	logger.Debug().Interface("stats", stats).Msgf("Reservation totals for last 24 hours: success=%d, failure=%d, pending=%d", success, failure, pending)

	for _, s := range stats.Usage28d {
		metrics.SetReservations28dCount(s.Result, s.Provider, s.Count)
	}

	return nil
}
