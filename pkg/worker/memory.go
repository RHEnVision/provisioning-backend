package worker

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

	if config.Telemetry.Enabled {
		job.TraceContext = make(map[string]string)
		otel.GetTextMapPropagator().Inject(ctx, job.TraceContext)
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

func (w *MemoryWorker) processJob(origCtx context.Context, job *Job) {
	if job == nil {
		zerolog.Ctx(origCtx).Error().Err(ErrJobNotFound).Msg("No job to process")
		return
	}

	ctx, logger, span := initJobContext(origCtx, job)
	defer span.End()

	if h, ok := w.handlers[job.Type]; ok {
		cCtx, cFunc := context.WithTimeout(ctx, config.Worker.Timeout)
		defer cFunc()
		h(cCtx, job)
	} else {
		span.SetStatus(codes.Error, "worker has not found handler for a job type")
		logger.Warn().Msgf("Memory worker handler not found for job type: %s", job.Type)
	}
}

func (w *MemoryWorker) Stats(_ context.Context) (Stats, error) {
	return Stats{}, nil
}
