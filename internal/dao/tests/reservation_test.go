//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createNoopReservation() *models.NoopReservation {
	return &models.NoopReservation{
		Reservation: models.Reservation{
			ID:        10,
			Provider:  models.ProviderTypeNoop,
			AccountID: 1,
			Status:    "Created",
		},
	}
}

func createInstancesReservation(reservationId int64) *models.ReservationInstance {
	return &models.ReservationInstance{
		ReservationID: reservationId,
		InstanceID:    "1",
	}
}

func createAWSReservation() *models.AWSReservation {
	return &models.AWSReservation{
		Reservation: models.Reservation{
			Provider:  models.ProviderTypeAWS,
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

func TestCreateNoop(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	err := reservationDao.CreateNoop(ctx, createNoopReservation())
	require.NoError(t, err)
	reservations, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, 1, len(reservations))
}

func TestCreateAWS(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	err := reservationDao.CreateAWS(ctx, createAWSReservation())
	require.NoError(t, err)
	reservations, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, 1, len(reservations))
}

func TestCreateInstance(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	reservation := createAWSReservation()
	err := reservationDao.CreateAWS(ctx, reservation)
	require.NoError(t, err)
	err = reservationDao.CreateInstance(ctx, createInstancesReservation(reservation.ID))
	require.NoError(t, err)
	reservations, err := reservationDao.ListInstances(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, 1, len(reservations))
}

func TestListReservation(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	awsReservation := createAWSReservation()
	err := reservationDao.CreateAWS(ctx, awsReservation)
	require.NoError(t, err)
	noopReservation := createNoopReservation()
	err = reservationDao.CreateNoop(ctx, noopReservation)
	require.NoError(t, err)
	reservations, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, 2, len(reservations))
}

func TestUpdateReservationIDForAWS(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	reservation := createAWSReservation()
	err := reservationDao.CreateAWS(ctx, reservation)
	require.NoError(t, err)
	var count int

	err = db.DB.Get(&count, "SELECT COUNT(*) FROM aws_reservation_details")
	require.NoError(t, err)

	assert.Equal(t, 1, count)

	err = reservationDao.UpdateReservationIDForAWS(ctx, reservation.ID, "2")
	require.NoError(t, err)
	var awsReservationId string
	err = db.DB.Get(&awsReservationId, "SELECT aws_reservation_id FROM aws_reservation_details WHERE reservation_id = $1", reservation.ID)
	require.NoError(t, err)

	assert.Equal(t, "2", awsReservationId)
}

func TestUpdateStatusReservation(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	err := reservationDao.CreateNoop(ctx, createNoopReservation())
	require.NoError(t, err)
	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	err = reservationDao.UpdateStatus(ctx, reservationsBefore[0].ID, "Edited", 0)
	require.NoError(t, err)

	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, "Edited", reservationsAfter[0].Status)
	assert.Equal(t, reservationsBefore[0].Step, reservationsAfter[0].Step)
}

func TestUpdateStepReservation(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	err := reservationDao.CreateNoop(ctx, createNoopReservation())
	require.NoError(t, err)

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	err = reservationDao.UpdateStatus(ctx, reservationsBefore[0].ID, "Edited", 42)
	require.NoError(t, err)

	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, reservationsBefore[0].Step+42, reservationsAfter[0].Step)
}

func TestDeleteReservation(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	noopReservation := createNoopReservation()
	err := reservationDao.CreateNoop(ctx, noopReservation)
	require.NoError(t, err)

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	err = reservationDao.Delete(ctx, reservationsBefore[0].ID)
	require.NoError(t, err)
	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, len(reservationsBefore)-1, len(reservationsAfter))
}

func TestFinishReservation(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	err := reservationDao.CreateNoop(ctx, createNoopReservation())
	require.NoError(t, err)

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	err = reservationDao.FinishWithSuccess(ctx, reservationsBefore[0].ID)
	require.NoError(t, err)

	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, reservationsBefore[0].ID, reservationsAfter[0].ID)
	assert.Equal(t, true, reservationsAfter[0].Success.Bool)
}

func TestFinishWithErrorReservation(t *testing.T) {
	reservationDao, ctx := setupReservation(t)
	defer teardownReservation(t)
	err := reservationDao.CreateNoop(ctx, createNoopReservation())
	require.NoError(t, err)

	reservationsBefore, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	err = reservationDao.FinishWithError(ctx, reservationsBefore[0].ID, "An error")
	require.NoError(t, err)

	reservationsAfter, err := reservationDao.List(ctx, 10, 0)
	require.NoError(t, err)

	assert.Equal(t, "An error", reservationsAfter[0].Error)
	assert.Equal(t, false, reservationsAfter[0].Success.Bool)
}
