// Import this package in all Go tests which need to enqueue a job. The implementation
// is to silently enqueue all incoming jobs without any effects. No jobs will be actually
// executed.
package stub

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

type enqueueCtxKeyType string

var (
	ContextReadError = errors.New("failed to find or convert dao stored in testing context")
	JobNotFoundError = errors.New("no job found")
)

var enqueueCtxKey enqueueCtxKeyType = "enqueuer-stub"

type hollowEnqueuer struct{}

type stubEnqueuer struct {
	enqueued []*worker.Job
}

func init() {
	queue.GetEnqueuer = getEnqueuer
}

// WithEnqueuer returns new context with Job enqueue struct that keeps the jobs
func WithEnqueuer(parent context.Context) context.Context {
	ctx := context.WithValue(parent, enqueueCtxKey, &stubEnqueuer{})
	return ctx
}

func EnqueuedJobs(ctx context.Context) []*worker.Job {
	enquer := getEnqueuerStub(ctx)
	return enquer.enqueued
}

func getEnqueuer(ctx context.Context) worker.JobEnqueuer {
	if enqueue := getEnqueuerStub(ctx); enqueue != nil {
		return enqueue
	}
	return hollowEnqueuer{}
}

func getEnqueuerStub(ctx context.Context) *stubEnqueuer {
	if enqueue, ok := ctx.Value(enqueueCtxKey).(*stubEnqueuer); ok {
		return enqueue
	}
	return nil
}

// Enqueue of hollow - default - enqueuer just ignores all enqueued jobs.
func (h hollowEnqueuer) Enqueue(_ context.Context, _ *worker.Job) error {
	return nil
}

func (s *stubEnqueuer) Enqueue(ctx context.Context, job *worker.Job) error {
	if job == nil {
		return fmt.Errorf("failed to enqueue: %w", JobNotFoundError)
	}
	s.enqueued = append(s.enqueued, job)
	return nil
}
