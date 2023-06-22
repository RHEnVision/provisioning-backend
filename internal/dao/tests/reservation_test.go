//go:build integration
// +build integration

package tests

import (
	"context"
	"math"
	"testing"
	"time"

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
			ID:         10,
			Steps:      1,
			StepTitles: []string{"Test step"},
			Provider:   models.ProviderTypeNoop,
			AccountID:  1,
			Status:     "Created",
		},
	}
}

func newReservationInstance(reservationId int64) *models.ReservationInstance {
	return &models.ReservationInstance{
		ReservationID: reservationId,
		InstanceID:    "1",
		Detail: models.ReservationInstanceDetail{
			PublicIPv4: "198.51.100.1",
		},
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
	ctx := identity.WithTenant(t, context.Background())
	reservationDao := dao.GetReservationDao(ctx)
	return reservationDao, ctx
}

func setupReservationOrg2(t *testing.T) (dao.ReservationDao, context.Context) {
	ctx := identity.WithTenantOrgId(t, context.Background(), "2")
	reservationDao := dao.GetReservationDao(ctx)
	return reservationDao, ctx
}

func TestReservationCreateNoop(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
		assert.Equal(t, res.AccountID, newRes.AccountID)
		assert.Equal(t, res.StepTitles, newRes.StepTitles)
		assert.Equal(t, res.Steps, newRes.Steps)
		assert.Equal(t, res.Status, newRes.Status)
		assert.Equal(t, time.Now().Year(), res.CreatedAt.Year())
	})
}

func TestReservationGetById(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		res := newNoopReservation()
		err := reservationDao.CreateNoop(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
		assert.Equal(t, res.AccountID, newRes.AccountID)
		assert.Equal(t, res.StepTitles, newRes.StepTitles)
		assert.Equal(t, res.Steps, newRes.Steps)
		assert.Equal(t, res.Status, newRes.Status)
		assert.Equal(t, time.Now().Year(), res.CreatedAt.Year())
	})

	t.Run("no rows", func(t *testing.T) {
		_, err := reservationDao.GetById(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrNoRows)
	})
}

func TestReservationCreateAWS(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		res := newAWSReservation()
		err := reservationDao.CreateAWS(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
		assert.Equal(t, res.AccountID, newRes.AccountID)
		assert.Equal(t, res.StepTitles, newRes.StepTitles)
		assert.Equal(t, res.Steps, newRes.Steps)
		assert.Equal(t, res.Status, newRes.Status)
		assert.Equal(t, time.Now().Year(), res.CreatedAt.Year())
	})
}

func TestReservationCreateGCP(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		res := newGCPReservation()
		err := reservationDao.CreateGCP(ctx, res)
		require.NoError(t, err)

		newRes, err := reservationDao.GetById(ctx, res.ID)
		require.NoError(t, err)
		assert.Equal(t, res.ID, newRes.ID)
		assert.Equal(t, res.AccountID, newRes.AccountID)
		assert.Equal(t, res.StepTitles, newRes.StepTitles)
		assert.Equal(t, res.Steps, newRes.Steps)
		assert.Equal(t, res.Status, newRes.Status)
		assert.Equal(t, time.Now().Year(), res.CreatedAt.Year())
	})
}

func TestReservationCreateAWSInstance(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		reservation := newAWSReservation()
		err := reservationDao.CreateAWS(ctx, reservation)
		require.NoError(t, err)

		instance := newReservationInstance(reservation.ID)
		err = reservationDao.CreateInstance(ctx, instance)
		require.NoError(t, err)

		instancesList, err := reservationDao.ListInstances(ctx, reservation.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(instancesList))
		assert.Equal(t, instance.InstanceID, instancesList[0].InstanceID)
		assert.Equal(t, instance.Detail.PublicIPv4, instancesList[0].Detail.PublicIPv4)
	})
}

func TestReservationList(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

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

func TestUnscopedUpdateAWSDetail(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

	t.Run("updates detail field", func(t *testing.T) {
		reservation := newAWSReservation()
		reservation.Detail = &models.AWSDetail{
			Region: "us-east-1",
			Amount: 1,
		}
		err := reservationDao.CreateAWS(ctx, reservation)
		require.NoError(t, err)
		newDetail := &models.AWSDetail{
			Region:     "us-east-1",
			Amount:     1,
			PubkeyName: "AWS name",
		}
		err = reservationDao.UnscopedUpdateAWSDetail(ctx, reservation.ID, newDetail)
		require.NoError(t, err, "failed to update reservation")

		reservationAfter, err := reservationDao.GetAWSById(ctx, reservation.ID)
		require.NoError(t, err)
		assert.Equal(t, "AWS name", reservationAfter.Detail.PubkeyName)
	})
}

func TestReservationUpdateIDForAWS(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer reset()

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
	defer reset()

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
	defer reset()

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
	defer reset()

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

func TestReservationRate(t *testing.T) {
	rdao, ctx := setupReservation(t)
	t.Run("allows slow reservations", func(t *testing.T) {
		defer reset()
		for i := 1; i <= 5; i++ {
			println(i)
			res := newNoopReservation()
			err := rdao.CreateNoop(ctx, res)
			require.NoError(t, err)
		}

		res := newNoopReservation()
		err := rdao.CreateNoop(ctx, res)
		require.ErrorContains(t, err, "failed to create reservation record")
	})

	rdao2, ctx2 := setupReservationOrg2(t)
	t.Run("throttles too fast reservations", func(t *testing.T) {
		defer reset()
		for i := 1; i <= 5; i++ {
			res := newNoopReservation()
			err := rdao.CreateNoop(ctx, res)
			require.NoError(t, err)
		}

		res := newNoopReservation()
		err := rdao2.CreateNoop(ctx2, res)
		require.NoError(t, err)
	})
}
