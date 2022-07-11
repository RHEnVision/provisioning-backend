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
)

// Queue is the main job queue
var Queue dejq.Jobs

const (
	TypeNoop            = "no_operation"
	TypePubkeyUploadAws = "pubkey_upload_aws"
)

func RegisterJobs(logger *zerolog.Logger) {
	logger.Debug().Msg("Initializing job queue")
	Queue.RegisterHandler(TypeNoop, HandleNoop)
	Queue.RegisterHandler(TypePubkeyUploadAws, HandlePubkeyUploadAWS)
}

func Initialize(ctx context.Context, logger *zerolog.Logger) error {
	var err error
	if config.Worker.Queue == "memory" {
		Queue, err = mem.NewClient(ctx, zerologr.New(logger))
	} else if config.Worker.Queue == "postgres" {
		Queue, err = postgres.NewClient(ctx, zerologr.New(logger), db.DB.DB,
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

func StartDequeueLoop(ctx context.Context, logger *zerolog.Logger) {
	logger.Debug().Msg("Starting dequeue loop")
	ctx = ctxval.WithLogger(ctx, logger)
	Queue.DequeueLoop(ctx)
}
