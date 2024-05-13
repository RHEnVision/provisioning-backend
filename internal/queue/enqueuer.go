package queue

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

var GetEnqueuer = func(ctx context.Context) worker.JobEnqueuer {
	panic("enqueuer not initialized")
}
