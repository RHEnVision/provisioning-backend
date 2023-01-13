// Import this package in all Go tests which need to enqueue a job. The implementation
// is to silently enqueue all incoming jobs without any effects. No jobs will be actually
// executed.
package stub

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

var (
	enqueuer worker.JobEnqueuer
	workers  worker.JobWorker
)

func getEnqueuer() worker.JobEnqueuer {
	return enqueuer
}

func init() {
	InitializeStub(context.Background())
	queue.GetEnqueuer = getEnqueuer
}

func InitializeStub(ctx context.Context) {
	wk := worker.NewMemoryClient()
	enqueuer = wk
	workers = wk
	workers.DequeueLoop(ctx)
}
