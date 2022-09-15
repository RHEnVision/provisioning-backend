package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type reservationDaoStub struct {
	lastId int64
	store  []*models.AWSReservation
}

func init() {
	dao.GetReservationDao = getReservationDao
}

func ReservationStubCount(ctx context.Context) (int, error) {
	pkdao, err := getReservationDaoStub(ctx)
	if err != nil {
		return 0, err
	}
	return len(pkdao.store), nil
}

func getReservationDao(ctx context.Context) (dao.ReservationDao, error) {
	return getReservationDaoStub(ctx)
}

func (stub *reservationDaoStub) NameForError() string {
	return "reservation"
}

func (stub *reservationDaoStub) CreateAWS(ctx context.Context, reservation *models.AWSReservation) error {
	reservation.ID = stub.lastId + 1
	stub.store = append(stub.store, reservation)
	stub.lastId++
	return nil
}

func (stub *reservationDaoStub) CreateGCP(ctx context.Context, reservation *models.GCPReservation) error {
	return nil
}

func (stub *reservationDaoStub) CreateNoop(ctx context.Context, reservation *models.NoopReservation) error {
	return nil
}

func (stub *reservationDaoStub) CreateInstance(ctx context.Context, reservation *models.ReservationInstance) error {
	return nil
}

func (stub *reservationDaoStub) GetById(ctx context.Context, id int64) (*models.Reservation, error) {
	return nil, nil
}

func (stub *reservationDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error) {
	return nil, nil
}

func (stub *reservationDaoStub) ListInstances(ctx context.Context, limit, offset int64) ([]*models.ReservationInstance, error) {
	return nil, nil
}

func (stub *reservationDaoStub) UpdateStatus(ctx context.Context, id int64, status string, addSteps int32) error {
	return nil
}

func (stub *reservationDaoStub) UpdateReservationIDForAWS(ctx context.Context, id int64, awsReservationId string) error {
	return nil
}

func (stub *reservationDaoStub) UpdateOperationNameForGCP(ctx context.Context, id int64, gcpOperationName string) error {
	return nil
}

func (stub *reservationDaoStub) FinishWithSuccess(ctx context.Context, id int64) error {
	return nil
}

func (stub *reservationDaoStub) FinishWithError(ctx context.Context, id int64, errorString string) error {
	return nil
}

func (stub *reservationDaoStub) Delete(ctx context.Context, id int64) error {
	return nil
}
