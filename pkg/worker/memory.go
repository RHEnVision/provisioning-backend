package worker

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
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
	if job == nil {
		return fmt.Errorf("unable to enqueue job: %w", ErrJobNotFound)
	}

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
	zerolog.Ctx(ctx).Info().Msg("Starting memory dequeuer")
	go w.dequeueLoop(ctx)
}

func (w *MemoryWorker) dequeueLoop(ctx context.Context) {
	for job := range w.todo {
		w.processJob(ctx, job)
	}
}

func (w *MemoryWorker) processJob(ctx context.Context, job *Job) {
	if job == nil {
		zerolog.Ctx(ctx).Error().Err(ErrJobNotFound).Msg("No job to process")
		return
	}

	if h, ok := w.handlers[job.Type]; ok {
		ctx, _ = contextLogger(ctx, job)
		cCtx, cFunc := context.WithTimeout(ctx, config.Worker.Timeout)
		defer cFunc()
		h(cCtx, job)
	} else {
		zerolog.Ctx(ctx).Warn().Msgf("Memory worker handler not found for job type: %s", job.Type)
	}
}

func (w *MemoryWorker) Stats(_ context.Context) (Stats, error) {
	return Stats{}, nil
}
