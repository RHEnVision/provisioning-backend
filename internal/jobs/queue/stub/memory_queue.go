// Import this package in all Go tests which need to enqueue a job. The implementation
// is to silently handle all incoming jobs without any effects.
package stub

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/go-logr/zerologr"
	"github.com/lzap/dejq"
	"github.com/lzap/dejq/mem"
	"github.com/rs/zerolog/log"
)

// memQueue is the main job stub memQueue
var memQueue dejq.Jobs

func getEnqueuer() queue.Enqueuer {
	return memQueue
}

func init() {
	InitializeStub()
	queue.GetEnqueuer = getEnqueuer
}

func InitializeStub() {
	var err error
	ctx := context.Background()
	memQueue, err = mem.NewClient(ctx, zerologr.New(&log.Logger))
	if err != nil {
		panic(err)
	}
	memQueue.DequeueLoop(ctx)
}
