package taskq

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/rs/zerolog"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
	"github.com/vmihailenco/taskq/v3/redisq"
)

func Initialize(_ context.Context, logger *zerolog.Logger) {
	logger.Debug().Msgf("Initializing '%s' job queue factory", config.Worker.Queue)

	switch config.Worker.Queue {
	case "memory":
		factory := memqueue.NewFactory()
		queue.JobQueue = factory.RegisterQueue(&taskq.QueueOptions{
			Name:    config.QueueName(),
			Storage: taskq.NewLocalStorage(),
		})
	case "redis":
		factory := redisq.NewFactory()
		queue.JobQueue = factory.RegisterQueue(&taskq.QueueOptions{
			Name:        config.QueueName(),
			Redis:       config.RedisDB(),
			WorkerLimit: config.Worker.Concurrency,
		})

	default:
		panic("unknown WORKER_QUEUE setting, expected values: memory, redis")
	}

	//taskq.SetLogger()
}
