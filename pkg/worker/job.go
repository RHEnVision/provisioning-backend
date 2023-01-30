package worker

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func init() {
	// makes UUID generation faster
	uuid.EnableRandPool()
}

type JobType string

type JobHandler func(ctx context.Context, job *Job)

type Job struct {
	// Random UUID for logging and tracing. It is generated randomly by Enqueue function when blank.
	ID uuid.UUID

	// Job type or "queue".
	Type JobType

	// Associated identity
	Identity identity.XRHID

	// Job arguments.
	Args any
}

var HandlerNotFoundErr = errors.New("handler not registered")

// JobEnqueuer sends Job messages into worker queue.
type JobEnqueuer interface {
	// Enqueue delivers a job to one of the backend workers.
	Enqueue(context.Context, *Job) error
}

// JobWorker receives and handles Job messages.
type JobWorker interface {
	// RegisterHandler registers an event listener for a particular type with an associated handler.
	RegisterHandler(JobType, JobHandler, any)

	// DequeueLoop starts one or more goroutines to dispatch incoming jobs.
	DequeueLoop(ctx context.Context)

	// Stop let's background workers to finish all jobs and terminates them. It blocks until workers are done.
	Stop(ctx context.Context)

	// Stats returns statistics. Not all implementations supports stats, some may return zero values.
	Stats(ctx context.Context) (Stats, error)
}

// Stats provides monitoring statistics.
type Stats struct {
	// Number of jobs currently in the queue
	EnqueuedJobs uint64

	// Number of jobs currently being processed
	InFlight int64
}
