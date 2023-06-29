package background

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/rs/zerolog"
)

// jobQueueMetricLoop is a background function that runs for all workers.
// It polls job queue statistics from Redis as well as in-flight counters.
//
//nolint:gosec
func jobQueueMetricLoop(ctx context.Context, sleep time.Duration, name string) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msgf("Started Redis statistics routine with tick interval %.2f seconds", sleep.Seconds())
	defer func() {
		logger.Debug().Msgf("Redis statistics routine exited")
	}()
	ticker := time.NewTicker(sleep)

	for {
		select {
		case <-ticker.C:
			stats := jq.Stats(ctx)
			logger.Debug().Msgf("Job queue statistics: enqueued=%d, in-flight=%d", stats.EnqueuedJobs, stats.InFlight)
			metrics.SetJobQueueSize(stats.EnqueuedJobs)
			metrics.SetJobQueueInFlight(name, stats.InFlight)

		case <-ctx.Done():
			ticker.Stop()
			logger.Debug().Msg("Stopping job queue metric loop")
			return
		}
	}
}
