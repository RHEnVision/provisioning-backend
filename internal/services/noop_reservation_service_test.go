package services

import (
	"context"
	"testing"

	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"

	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/lzap/dejq"
)

func TestEnqueueNoopJob(t *testing.T) {
	pj := dejq.PendingJob{
		Type: jobs.TypeNoop,
		Body: &jobs.NoopJobArgs{
			AccountID:     1,
			ReservationID: 0,
		},
	}
	err := jobs.Enqueue(context.Background(), pj)
	if err != nil {
		panic(err)
	}
}

func TestCreateNoopReservation(t *testing.T) {
	// TODO full service test
}
