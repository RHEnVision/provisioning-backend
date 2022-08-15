package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/go-logr/zerologr"
	"github.com/lzap/dejq"
	"github.com/lzap/dejq/mem"
	"github.com/lzap/dejq/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// queue is the main job queue
var queue dejq.Jobs

const (
	TypeNoop              = "no_operation"
	TypePubkeyUploadAws   = "pubkey_upload_aws"
	TypeLaunchInstanceAws = "launch_instance_aws"
)

func RegisterJobs(logger *zerolog.Logger) {
	logger.Debug().Msg("Initializing job queue")
	queue.RegisterHandler(TypeNoop, HandleNoop)
	queue.RegisterHandler(TypePubkeyUploadAws, HandlePubkeyUploadAWS)
	queue.RegisterHandler(TypeLaunchInstanceAws, HandleLaunchInstanceAWS)
}

func Initialize(ctx context.Context, logger *zerolog.Logger) error {
	var err error
	if config.Worker.Queue == "memory" {
		queue, err = mem.NewClient(ctx, zerologr.New(logger))
	} else if config.Worker.Queue == "postgres" {
		queue, err = postgres.NewClient(ctx, zerologr.New(logger), db.DB.DB,
			config.Worker.Concurrency,
			time.Duration(config.Worker.HeartbeatSec)*time.Second,
			config.Worker.MaxBeats)
	} else if config.Worker.Queue == "sqs" {
		panic("SQS queue implementation is not supported")
	}
	if err != nil {
		return fmt.Errorf("cannot initialize queue: %w", err)
	}
	return nil
}

func Enqueue(ctx context.Context, jobs ...dejq.PendingJob) error {
	if queue == nil {
		panic("job queue was not initialized yet, for tests import internal/jobs/stub")
	}
	err := queue.Enqueue(ctx, jobs...)
	if err != nil {
		return fmt.Errorf("enqueue error: %w", err)
	}
	return nil
}

func StartDequeueLoop(ctx context.Context, logger *zerolog.Logger) {
	logger.Debug().Msg("Starting dequeue loop")
	ctx = ctxval.WithLogger(ctx, logger)
	queue.DequeueLoop(ctx)
}

func StopDequeueLoop() {
	queue.Stop()
}

func InitializeStub() {
	var err error
	ctx := context.Background()
	queue, err = mem.NewClient(ctx, zerologr.New(&log.Logger))
	if err != nil {
		panic(err)
	}
	queue.DequeueLoop(ctx)
}
