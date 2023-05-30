package worker

import (
	"context"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/rs/zerolog"

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

	// Associated account.
	AccountID int64

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

func (jt JobType) String() string {
	return string(jt)
}

// Stats provides monitoring statistics.
type Stats struct {
	// Number of jobs currently in the queue. This is a global value - all clients see the same value.
	EnqueuedJobs uint64

	// Number of jobs currently being processed. Local value - each client has its own number.
	InFlight int64
}

func contextLogger(ctx context.Context, job *Job) context.Context {
	accountId := job.AccountID
	id := job.Identity
	logger := zerolog.Ctx(ctx).With().
		Str("job_id", job.ID.String()).
		Int64("account_id", accountId).
		Str("account_number", id.Identity.AccountNumber).
		Str("org_id", id.Identity.OrgID).Logger()
	newContext := logger.WithContext(ctx)
	newContext = ctxval.WithIdentity(newContext, id)
	newContext = ctxval.WithAccountId(newContext, accountId)

	return newContext
}
