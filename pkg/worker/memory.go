package worker

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/google/uuid"
)

type MemoryWorker struct {
	handlers map[JobType]JobHandler
	todo     chan *Job
}

func NewMemoryClient() *MemoryWorker {
	return &MemoryWorker{
		handlers: make(map[JobType]JobHandler),
		todo:     make(chan *Job),
	}
}

func (w *MemoryWorker) RegisterHandler(jtype JobType, handler JobHandler, _ any) {
	w.handlers[jtype] = handler
}

func (w *MemoryWorker) Enqueue(ctx context.Context, job *Job) error {
	var err error

	if job.ID == uuid.Nil {
		job.ID, err = uuid.NewRandom()
		if err != nil {
			return fmt.Errorf("unable to generate UUID: %w", err)
		}
	}

	w.todo <- job
	return nil
}

func (w *MemoryWorker) Stop(_ context.Context) {
	close(w.todo)
}

func (w *MemoryWorker) DequeueLoop(ctx context.Context) {
	go w.dequeueLoop(ctx)
}

func (w *MemoryWorker) dequeueLoop(ctx context.Context) {
	for job := range w.todo {
		if h, ok := w.handlers[job.Type]; ok {
			h(ctx, job)
		} else {
			ctxval.Logger(ctx).Warn().Msgf("Memory worker handler not found for job type: %s", job.Type)
		}
	}
}

func (w *MemoryWorker) Stats(_ context.Context) (Stats, error) {
	return Stats{}, nil
}
