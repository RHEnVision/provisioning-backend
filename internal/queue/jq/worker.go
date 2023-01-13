package jq

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/rs/zerolog"
)

var (
	enqueuer worker.JobEnqueuer
	workers  worker.JobWorker
)

func getEnqueuer() worker.JobEnqueuer {
	return enqueuer
}

func init() {
	queue.GetEnqueuer = getEnqueuer
}

func RegisterJobs(logger *zerolog.Logger) {
	logger.Debug().Msg("Registering job queue handlers and interfaces")
	workers.RegisterHandler(jobs.TypeNoop, jobs.HandleNoop, jobs.NoopJobArgs{})
	workers.RegisterHandler(jobs.TypeLaunchInstanceAws, jobs.HandleLaunchInstanceAWS, jobs.LaunchInstanceAWSTaskArgs{})
	workers.RegisterHandler(jobs.TypeLaunchInstanceGcp, jobs.HandleLaunchInstanceGCP, jobs.LaunchInstanceGCPTaskArgs{})
}

func Initialize(_ context.Context, logger *zerolog.Logger) error {
	logger.Debug().Msgf("Initializing '%s' job queue", config.Worker.Queue)

	switch config.Worker.Queue {
	case "memory":
		wk := worker.NewMemoryClient()
		enqueuer = wk
		workers = wk
	case "redis":
		wk, err := worker.NewRedisWorker(config.RedisHostAndPort(),
			config.Application.Cache.Redis.User, config.Application.Cache.Redis.Password,
			config.Application.Cache.Redis.DB, "provisioning-job-queue",
			config.Worker.PollInterval, config.Worker.MaxThreads)
		if err != nil {
			return fmt.Errorf("cannot initialize redis worker queue: %w", err)
		}
		enqueuer = wk
		workers = wk
	default:
		panic("unknown WORKER_QUEUE setting, expected values: memory, redis, postgres")
	}

	return nil
}

func StartDequeueLoop(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Starting dequeue loop")
	workers.DequeueLoop(ctx)
}

func StopDequeueLoop(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Stopping dequeue loop")
	workers.Stop(ctx)
}

func Stats(ctx context.Context) worker.Stats {
	stats, err := workers.Stats(ctx)
	if err != nil {
		ctxval.Logger(ctx).Error().Err(err).Msg("Unable to get queue stats")
		return worker.Stats{}
	}

	return stats
}
