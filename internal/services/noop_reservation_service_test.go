package services_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"

	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/lzap/dejq"
)

func TestEnqueueNoopJob(t *testing.T) {
	pj := dejq.PendingJob{
		Type: queue.TypeNoop,
		Body: &jobs.NoopJobArgs{
			AccountID:     1,
			ReservationID: 0,
		},
	}
	err := queue.GetEnqueuer().Enqueue(context.Background(), pj)
	if err != nil {
		panic(err)
	}
}

func TestCreateNoopReservationHandler(t *testing.T) {
	// TODO full service test
}
