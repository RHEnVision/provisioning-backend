package background

import (
	"context"
	"math/rand"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
)

// jobQueueMetricLoop is a background function that runs for all workers.
// It polls job queue statistics from Redis as well as in-flight counters.
// It waits a random delay between 0 and sleep for better polling spread
// since job queue size is a global metric (redis queue length).
//
//nolint:gosec
func jobQueueMetricLoop(ctx context.Context, sleep time.Duration, name string) {
	logger := ctxval.Logger(ctx)

	// spread polling intervals
	randSleep := rand.Int63() % sleep.Milliseconds()
	logger.Debug().Msgf("Job queue metric delay %dms", randSleep)
	time.Sleep(time.Duration(randSleep) * time.Millisecond)
	ticker := time.NewTicker(sleep)

	for {
		select {
		case <-ticker.C:
			stats := jq.Stats(ctx)
			metrics.SetJobQueueSize(stats.EnqueuedJobs)
			metrics.SetJobQueueInFlight(name, stats.InFlight)

		case <-ctx.Done():
			ticker.Stop()
			logger.Debug().Msg("Stopping job queue metric loop")
			return
		}
	}
}
