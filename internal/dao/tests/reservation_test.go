//go:build integration
// +build integration

package tests

import (
	"context"
	"math"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newNoopReservation() *models.NoopReservation {
	return &models.NoopReservation{
		Reservation: models.Reservation{
			ID:        10,
			Provider:  models.ProviderTypeNoop,
			AccountID: 1,
			Status:    "Created",
		},
	}
}

func newInstancesReservation(reservationId int64) *models.ReservationInstance {
	return &models.ReservationInstance{
		ReservationID: reservationId,
		InstanceID:    "1",
	}
}

func newAWSReservation() *models.AWSReservation {
	return &models.AWSReservation{
		Reservation: models.Reservation{
			Provider:  models.ProviderTypeAWS,
			AccountID: 1,
			Status:    "Created",
		},
		PubkeyID: 1,
	}
}

func newGCPReservation() *models.GCPReservation {
	return &models.GCPReservation{
		Reservation: models.Reservation{
			Provider:  models.ProviderTypeGCP,
			AccountID: 1,
			Status:    "Created",
		},
		PubkeyID: 1,
	}
}

func setupReservation(t *testing.T) (dao.ReservationDao, context.Context) {
	setup()
	ctx := identity.WithTenant(t, context.Background())
	reservationDao, err := dao.GetReservationDao(ctx)
	require.NoError(t, err)
	return reservationDao, ctx
}

func teardownReservation(_ *testing.T) {
	teardown()
}

func TestReservationCreateNoop(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
	})
}

func TestReservationGetById(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
	})

	t.Run("no rows", func(t *testing.T) {
		_, err := reservationDao.GetById(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrNoRows)
	})
}

func TestReservationCreateAWS(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		res := newAWSReservation()
		err := reservationDao.CreateAWS(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
	})
}

func TestReservationCreateGCP(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		res := newGCPReservation()
		err := reservationDao.CreateGCP(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
	})
}

func TestReservationCreateAWSInstance(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		reservation := newAWSReservation()
		err := reservationDao.CreateAWS(ctx, reservation)
		require.NoError(t, err)

		err = reservationDao.CreateInstance(ctx, newInstancesReservation(reservation.ID))
		require.NoError(t, err)

		reservations, err := reservationDao.ListInstances(ctx, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(reservations))
	})
}

func TestReservationList(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("empty", func(t *testing.T) {
		reservations, err := reservationDao.List(ctx, 10, 0)
		require.NoError(t, err)
		require.Empty(t, reservations)
	})

	t.Run("success", func(t *testing.T) {
		awsReservation := newAWSReservation()
		err := reservationDao.CreateAWS(ctx, awsReservation)
		require.NoError(t, err)

		noopReservation := newNoopReservation()
		err = reservationDao.CreateNoop(ctx, noopReservation)
		require.NoError(t, err)

		reservations, err := reservationDao.List(ctx, 10, 0)
		require.NoError(t, err)
		assert.Equal(t, 2, len(reservations))
	})
}

func TestReservationUpdateIDForAWS(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		reservation := newAWSReservation()
		err := reservationDao.CreateAWS(ctx, reservation)
		require.NoError(t, err)
		var count int

		err = db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM aws_reservation_details").Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		err = reservationDao.UpdateReservationIDForAWS(ctx, reservation.ID, "r-8954738954")
		require.NoError(t, err)

		var awsReservationId string
		query := "SELECT aws_reservation_id FROM aws_reservation_details WHERE reservation_id = $1"
		err = db.Pool.QueryRow(ctx, query, reservation.ID).Scan(&awsReservationId)
		require.NoError(t, err)

		assert.Equal(t, "r-8954738954", awsReservationId)
	})
}

func TestReservationUpdateStatus(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("status text", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		err = reservationDao.UpdateStatus(ctx, res.ID, "Edited", 0)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, "Edited", newRes.Status)
		assert.Equal(t, res.Step, newRes.Step)
	})

	t.Run("step", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		err = reservationDao.UpdateStatus(ctx, res.ID, "New step", 1)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.Step+1, newRes.Step)
	})
}

func TestReservationDelete(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		err = reservationDao.Delete(ctx, res.ID)
		require.NoError(t, err)

		_, err = reservationDao.GetById(ctx, res.ID)
		require.ErrorIs(t, err, dao.ErrNoRows)
	})

	t.Run("mismatch", func(t *testing.T) {
		err := reservationDao.Delete(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})
}

func TestReservationFinish(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)

	t.Run("with success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		err = reservationDao.FinishWithSuccess(ctx, res.ID)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.True(t, newRes.Success.Valid)
		assert.True(t, newRes.Success.Bool)
	})

	t.Run("with error", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		err = reservationDao.FinishWithError(ctx, res.ID, "error")
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.True(t, newRes.Success.Valid)
		assert.False(t, newRes.Success.Bool)
		assert.Equal(t, "error", newRes.Error)
	})

	t.Run("mismatch success", func(t *testing.T) {
		err := reservationDao.FinishWithSuccess(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})

	t.Run("mismatch error", func(t *testing.T) {
		err := reservationDao.FinishWithError(ctx, math.MaxInt64, "")
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})
}
