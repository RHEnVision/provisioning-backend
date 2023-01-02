// Import this package in all Go tests which need to enqueue a job. Tasks are enqueued, however,
// no tasks are being picked up. This is only useful in unit tests.
package stub

import (
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
)

func init() {
	factory := memqueue.NewFactory()

	// create stub queue
	queue.JobQueue = factory.RegisterQueue(&taskq.QueueOptions{
		Name:        "stub-queue",
		WorkerLimit: 1,
	})

	// override the registration function to create noop handlers
	queue.RegisterTask = func(_ string, _ any) *taskq.Task {
		return taskq.RegisterTask(&taskq.TaskOptions{
			Name: "stub",
			Handler: func(_ *taskq.Message) error {
				return nil
			},
		})
	}

	// unregister already registered handlers
	taskq.Tasks.Reset()
}
