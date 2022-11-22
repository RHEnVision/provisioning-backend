package dejq

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/go-logr/zerologr"
	"github.com/lzap/dejq"
	"github.com/lzap/dejq/mem"
	"github.com/lzap/dejq/postgres"
	"github.com/lzap/dejq/redis"
	"github.com/rs/zerolog"
)

// dejqQueue is the main job dejqQueue
var dejqQueue dejq.Jobs

func getEnqueuer() queue.Enqueuer {
	return dejqQueue
}

func init() {
	queue.GetEnqueuer = getEnqueuer
}

func RegisterJobs(logger *zerolog.Logger) {
	logger.Debug().Msg("Registering job queue handlers")
	dejqQueue.RegisterHandler(queue.TypeNoop, jobs.HandleNoop)
	dejqQueue.RegisterHandler(queue.TypeEnsurePubkeyOnAws, jobs.HandleEnsurePubkeyOnAWS)
	dejqQueue.RegisterHandler(queue.TypeLaunchInstanceAws, jobs.HandleLaunchInstanceAWS)
	dejqQueue.RegisterHandler(queue.TypeLaunchInstanceGcp, jobs.HandleLaunchInstanceGCP)
}

func Initialize(ctx context.Context, logger *zerolog.Logger) error {
	logger.Debug().Msgf("Initializing '%s' job queue", config.Worker.Queue)

	var err error
	switch config.Worker.Queue {
	case "memory":
		dejqQueue, err = mem.NewClient(ctx, zerologr.New(logger))
	case "redis":
		dejqQueue, err = redis.NewClient(ctx, zerologr.New(logger), config.RedisHostAndPort(),
			config.Application.Cache.Redis.User, config.Application.Cache.Redis.Password,
			config.Application.Cache.Redis.DB, "provisioning-job-queue")
	case "postgres":
		// TODO dejq must be refactored to use PGX too
		dejqQueue, err = postgres.NewClient(ctx, zerologr.New(logger), nil,
			config.Worker.Concurrency,
			config.Worker.Heartbeat,
			config.Worker.MaxBeats)
	default:
		panic("unknown WORKER_QUEUE setting, expected values: memory, redis, postgres")
	}

	if err != nil {
		return fmt.Errorf("cannot initialize dejqQueue: %w", err)
	}

	if dejqQueue == nil {
		panic("dejq queue was not initialized, check configuration")
	}

	return nil
}

func StartDequeueLoop(ctx context.Context, logger *zerolog.Logger) {
	loggerWithQueue := logger.With().Str("queue_type", config.Worker.Queue).Logger()
	loggerWithQueue.Debug().Msg("Starting dequeue loop")
	ctx = ctxval.WithLogger(ctx, &loggerWithQueue)
	dejqQueue.DequeueLoop(ctx)
}

func StopDequeueLoop() {
	dejqQueue.Stop()
}
