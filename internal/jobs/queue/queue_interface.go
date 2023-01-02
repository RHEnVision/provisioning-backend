package queue

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/vmihailenco/taskq/v3"
)

var JobQueue taskq.Queue

// RegisterTask registers a handler. When called from unit tests, the implementation
// returns noop handler function.
var RegisterTask = func(name string, handler any) *taskq.Task {
	return taskq.RegisterTask(&taskq.TaskOptions{
		Name:    name,
		Handler: handler,
	})
}

func StartQueues(ctx context.Context, logger *zerolog.Logger) {
	// Queues are started automatically via RegisterQueue - no action needed
}

func StopQueues(logger *zerolog.Logger) {
	if err := JobQueue.Consumer().Stop(); err != nil {
		logger.Error().Err(err).Msg("Job dequeue loop error")
	}
}
