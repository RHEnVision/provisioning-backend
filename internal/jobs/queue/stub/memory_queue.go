// Import this package in all Go tests which need to enqueue a job. The implementation
// is to silently enqueue all incoming jobs without any effects. No jobs will be actually
// executed. Tests for job execution must use the dejq package (memory driver).
package stub

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/go-logr/zerologr"
	"github.com/lzap/dejq"
	"github.com/lzap/dejq/mem"
)

// memQueue is the main job stub memQueue
var memQueue dejq.Jobs

func getEnqueuer() queue.Enqueuer {
	return memQueue
}

func init() {
	InitializeStub(context.Background())
	queue.GetEnqueuer = getEnqueuer
}

func InitializeStub(ctx context.Context) {
	var err error

	logger := ctxval.Logger(ctx)
	memQueue, err = mem.NewClient(ctx, zerologr.New(logger))
	if err != nil {
		panic(err)
	}
	memQueue.DequeueLoop(ctx)
}
