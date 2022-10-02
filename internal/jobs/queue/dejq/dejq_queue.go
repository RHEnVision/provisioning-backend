package dejq

import (
	"context"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/go-logr/zerologr"
	"github.com/lzap/dejq"
	"github.com/lzap/dejq/mem"
	"github.com/lzap/dejq/postgres"
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
	logger.Debug().Msg("Initializing job queue")
	dejqQueue.RegisterHandler(queue.TypeNoop, jobs.HandleNoop)
	dejqQueue.RegisterHandler(queue.TypePubkeyUploadAws, jobs.HandlePubkeyUploadAWS)
	dejqQueue.RegisterHandler(queue.TypeLaunchInstanceAws, jobs.HandleLaunchInstanceAWS)
	dejqQueue.RegisterHandler(queue.TypeLaunchInstanceGcp, jobs.HandleLaunchInstanceGCP)
}

func Initialize(ctx context.Context, logger *zerolog.Logger) error {
	var err error
	if config.Worker.Queue == "memory" {
		dejqQueue, err = mem.NewClient(ctx, zerologr.New(logger))
	} else if config.Worker.Queue == "postgres" {
		// TODO dejq must be refactored to use PGX too
		dejqQueue, err = postgres.NewClient(ctx, zerologr.New(logger), nil,
			config.Worker.Concurrency,
			time.Duration(config.Worker.HeartbeatSec)*time.Second,
			config.Worker.MaxBeats)
	} else if config.Worker.Queue == "sqs" {
		panic("SQS dejqQueue implementation is not supported")
	}
	if err != nil {
		return fmt.Errorf("cannot initialize dejqQueue: %w", err)
	}
	return nil
}

func StartDequeueLoop(ctx context.Context, logger *zerolog.Logger) {
	logger.Debug().Msg("Starting dequeue loop")
	ctx = ctxval.WithLogger(ctx, logger)
	dejqQueue.DequeueLoop(ctx)
}

func StopDequeueLoop() {
	dejqQueue.Stop()
}
