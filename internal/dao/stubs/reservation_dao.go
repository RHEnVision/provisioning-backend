package stubs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type reservationDaoStub struct {
	storeAWS   []*models.AWSReservation
	storeAzure []*models.AzureReservation
	storeGCP   []*models.GCPReservation
	instances  map[int64][]*models.ReservationInstance
}

func init() {
	dao.GetReservationDao = getReservationDao
}

func AWSReservationStubCount(ctx context.Context) int {
	resDao := getReservationDaoStub(ctx)
	return len(resDao.storeAWS)
}

func AzureReservationStubCount(ctx context.Context) int {
	resDao := getReservationDaoStub(ctx)
	return len(resDao.storeAzure)
}

func GCPReservationStubCount(ctx context.Context) int {
	resDao := getReservationDaoStub(ctx)
	return len(resDao.storeGCP)
}

func getReservationDao(ctx context.Context) dao.ReservationDao {
	return getReservationDaoStub(ctx)
}

func (stub *reservationDaoStub) CreateAWS(ctx context.Context, reservation *models.AWSReservation) error {
	reservation.ID = int64(len(stub.storeAWS)) + 1
	stub.storeAWS = append(stub.storeAWS, reservation)
	return nil
}

func (stub *reservationDaoStub) CreateAzure(ctx context.Context, reservation *models.AzureReservation) error {
	reservation.ID = int64(len(stub.storeAzure)) + 1
	stub.storeAzure = append(stub.storeAzure, reservation)
	return nil
}

func (stub *reservationDaoStub) CreateGCP(ctx context.Context, reservation *models.GCPReservation) error {
	reservation.ID = int64(len(stub.storeGCP)) + 1
	stub.storeGCP = append(stub.storeGCP, reservation)
	return nil
}

func (stub *reservationDaoStub) CreateNoop(ctx context.Context, reservation *models.NoopReservation) error {
	return nil
}

func (stub *reservationDaoStub) CreateInstance(ctx context.Context, resInstance *models.ReservationInstance) error {
	resId := resInstance.ReservationID
	stub.instances[resId] = append(stub.instances[resId], resInstance)
	return nil
}

func (stub *reservationDaoStub) GetById(ctx context.Context, id int64) (*models.Reservation, error) {
	for _, awsReservation := range stub.storeAWS {
		if awsReservation.AccountID == ctxAccountId(ctx) && awsReservation.ID == id {
			return &awsReservation.Reservation, nil
		}
	}
	return nil, dao.ErrNoRows
}

func (stub *reservationDaoStub) GetAWSById(ctx context.Context, id int64) (*models.AWSReservation, error) {
	for _, awsReservation := range stub.storeAWS {
		if awsReservation.AccountID == ctxAccountId(ctx) && awsReservation.ID == id {
			return awsReservation, nil
		}
	}
	return nil, dao.ErrNoRows
}

func (stub *reservationDaoStub) GetAzureById(ctx context.Context, id int64) (*models.AzureReservation, error) {
	for _, azureReservation := range stub.storeAzure {
		if azureReservation.AccountID == ctxAccountId(ctx) && azureReservation.ID == id {
			return azureReservation, nil
		}
	}
	return nil, dao.ErrNoRows
}

func (stub *reservationDaoStub) GetGCPById(ctx context.Context, id int64) (*models.GCPReservation, error) {
	for _, gcpReservation := range stub.storeGCP {
		if gcpReservation.AccountID == ctxAccountId(ctx) && gcpReservation.ID == id {
			return gcpReservation, nil
		}
	}
	return nil, dao.ErrNoRows
}

func (stub *reservationDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error) {
	return nil, nil
}

func (stub *reservationDaoStub) ListInstances(ctx context.Context, reservationId int64) ([]*models.ReservationInstance, error) {
	return stub.instances[reservationId], nil
}

func (stub *reservationDaoStub) UpdateStatus(ctx context.Context, id int64, status string, addSteps int32) error {
	return nil
}

func (stub *reservationDaoStub) UnscopedUpdateAWSDetail(ctx context.Context, id int64, awsDetail *models.AWSDetail) error {
	res, err := stub.GetAWSById(ctx, id)
	if err != nil {
		return fmt.Errorf("stubbed lookup of AWS reservation failed: %w", err)
	}
	res.Detail = awsDetail
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

func (stub *reservationDaoStub) Cleanup(ctx context.Context) error {
	return nil
}

func (stub *reservationDaoStub) UpdateReservationInstance(ctx context.Context, reservationID int64, instance *clients.InstanceDescription) error {
	for _, instRes := range stub.instances[reservationID] {
		if instRes.InstanceID == instance.ID {
			instRes.Detail.PublicIPv4 = instance.PublicIPv4
		}
	}
	return nil
}
