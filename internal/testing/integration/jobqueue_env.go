package integration

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/rs/zerolog/log"
)

func InitJobQueueEnvironment(ctx context.Context) {
	err := jq.Initialize(ctx, &log.Logger)
	if err != nil {
		panic(fmt.Errorf("cannot initialize job queue: %w", err))
	}
	jq.RegisterJobs(&log.Logger)
	jq.StartDequeueLoop(ctx)
}

func CloseJobQueueEnvironment(ctx context.Context) {
	jq.StopDequeueLoop(ctx)
}
