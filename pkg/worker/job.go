package worker

import (
	"context"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
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
	Identity identity.Principal

	// For logging purposes
	TraceID string

	// For logging purposes
	EdgeID string

	// Job arguments.
	Args any
}

var ErrHandlerNotFound = errors.New("handler not registered")

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

func contextLogger(origCtx context.Context, job *Job) (context.Context, *zerolog.Logger) {
	if job == nil {
		zerolog.Ctx(origCtx).Warn().Err(ErrJobNotFound).Msg("No job, context not changed")
		return origCtx, nil
	}

	ctx := logging.WithJobId(origCtx, job.ID.String())
	ctx = identity.WithIdentity(ctx, job.Identity)
	ctx = logging.WithTraceId(ctx, job.TraceID)
	ctx = logging.WithEdgeRequestId(ctx, job.EdgeID)
	ctx = identity.WithAccountId(ctx, job.AccountID)
	ctx = logging.WithJobId(ctx, job.ID.String())
	ctx = logging.WithJobType(ctx, job.Type.String())

	logger := zerolog.Ctx(ctx)
	logger = ptr.To(logger.With().
		Int64("account_id", job.AccountID).
		Str("org_id", job.Identity.Identity.OrgID).
		Str("account_number", job.Identity.Identity.AccountNumber).
		Str("trace_id", job.TraceID).
		Str("request_id", job.EdgeID).
		Str("job_id", job.ID.String()).
		Str("job_type", job.Type.String()).
		Interface("job_args", job.Args).
		Logger())
	return logger.WithContext(ctx), logger
}
