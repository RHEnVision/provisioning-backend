//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getReservationDao(t *testing.T) (dao.ReservationDao, context.Context) {
	ctx := identity.WithTenant(t, context.Background())
	reservationDao := dao.GetReservationDao(ctx)
	return reservationDao, ctx
}

// read reservation until it is finished (max. 2 seconds)
func waitForReservation(t *testing.T, id int64) *models.Reservation {
	reservationDao, ctx := getReservationDao(t)

	var updatedRes *models.Reservation
	var err error
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		updatedRes, err = reservationDao.GetById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, updatedRes.ID)

		if updatedRes.Success.Valid {
			break
		}
	}
	return updatedRes
}

func TestRedisNoopSuccess(t *testing.T) {
	reservationDao, ctx := getReservationDao(t)
	defer reset()

	res := &models.NoopReservation{
		Reservation: models.Reservation{
			AccountID:  1,
			Steps:      1,
			StepTitles: []string{"Test step"},
			Provider:   models.ProviderTypeNoop,
			Status:     "Created",
		},
	}
	err := reservationDao.CreateNoop(ctx, res)
	require.NoError(t, err)
	require.NotZero(t, res.ID)

	job := worker.Job{
		AccountID: 1,
		Type:      jobs.TypeNoop,
		Args: jobs.NoopJobArgs{
			ReservationID: res.ID,
			Fail:          false,
		},
	}

	err = queue.GetEnqueuer(ctx).Enqueue(context.Background(), &job)
	require.NoError(t, err)
	updatedRes := waitForReservation(t, res.ID)

	require.True(t, updatedRes.Success.Valid)
	require.True(t, updatedRes.Success.Bool)
	require.Empty(t, updatedRes.Error)
	require.Equal(t, int32(1), updatedRes.Step)
	require.Equal(t, int32(1), updatedRes.Steps)
	require.Equal(t, "No operation finished", updatedRes.Status)
}

func TestRedisNoopFailure(t *testing.T) {
	reservationDao, ctx := getReservationDao(t)
	defer reset()

	res := &models.NoopReservation{
		Reservation: models.Reservation{
			AccountID:  1,
			Steps:      1,
			StepTitles: []string{"Test step"},
			Provider:   models.ProviderTypeNoop,
			Status:     "Created",
		},
	}
	err := reservationDao.CreateNoop(ctx, res)
	require.NoError(t, err)
	require.NotZero(t, res.ID)

	job := worker.Job{
		AccountID: 1,
		Type:      jobs.TypeNoop,
		Args: jobs.NoopJobArgs{
			ReservationID: res.ID,
			Fail:          true,
		},
	}

	err = queue.GetEnqueuer(ctx).Enqueue(context.Background(), &job)
	require.NoError(t, err)
	updatedRes := waitForReservation(t, res.ID)

	require.True(t, updatedRes.Success.Valid)
	require.False(t, updatedRes.Success.Bool)
	require.Equal(t, "job failed on request", updatedRes.Error)
	require.Equal(t, int32(1), updatedRes.Step)
	require.Equal(t, int32(1), updatedRes.Steps)
	require.Equal(t, "No operation finished", updatedRes.Status)
}

func TestRedisNoopTimeout(t *testing.T) {
	reservationDao, ctx := getReservationDao(t)
	defer reset()

	res := &models.NoopReservation{
		Reservation: models.Reservation{
			AccountID:  1,
			Steps:      1,
			StepTitles: []string{"Test step"},
			Provider:   models.ProviderTypeNoop,
			Status:     "Created",
		},
	}
	err := reservationDao.CreateNoop(ctx, res)
	require.NoError(t, err)
	require.NotZero(t, res.ID)

	job := worker.Job{
		AccountID: 1,
		Type:      jobs.TypeNoop,
		Args: jobs.NoopJobArgs{
			ReservationID: res.ID,
			Fail:          false,
			Sleep:         1300 * time.Millisecond,
		},
	}

	err = queue.GetEnqueuer(ctx).Enqueue(context.Background(), &job)
	require.NoError(t, err)
	updatedRes := waitForReservation(t, res.ID)

	require.True(t, updatedRes.Success.Valid)
	require.False(t, updatedRes.Success.Bool)
	require.Equal(t, "Timeout", updatedRes.Status)
}
