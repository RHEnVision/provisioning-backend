package dao

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

var GetAccountDao func(ctx context.Context) (AccountDao, error)

type AccountDao interface {
	Create(ctx context.Context, pk *models.Account) error
	GetById(ctx context.Context, id int64) (*models.Account, error)
	GetOrCreateByIdentity(ctx context.Context, orgId string, accountNumber string) (*models.Account, error)
	GetByOrgId(ctx context.Context, orgId string) (*models.Account, error)
	GetByAccountNumber(ctx context.Context, number string) (*models.Account, error)
	List(ctx context.Context, limit, offset int64) ([]*models.Account, error)
}

var GetPubkeyDao func(ctx context.Context) (PubkeyDao, error)

type PubkeyDao interface {
	Create(ctx context.Context, pk *models.Pubkey) error
	Update(ctx context.Context, pk *models.Pubkey) error
	GetById(ctx context.Context, id int64) (*models.Pubkey, error)
	List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error)
	Delete(ctx context.Context, id int64) error
}

var GetPubkeyResourceDao func(ctx context.Context) (PubkeyResourceDao, error)

type PubkeyResourceDao interface {
	GetResourceByProviderType(ctx context.Context, pubkeyId int64, provider models.ProviderType) (*models.PubkeyResource, error)
	ListByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error)
	Create(ctx context.Context, pkr *models.PubkeyResource) error
	Update(ctx context.Context, pkr *models.PubkeyResource) error
	Delete(ctx context.Context, id int64) error
}

var GetReservationDao func(ctx context.Context) (ReservationDao, error)

type ReservationDao interface {
	CreateNoop(ctx context.Context, reservation *models.NoopReservation) error
	CreateAWS(ctx context.Context, reservation *models.AWSReservation) error
	List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	Finish(ctx context.Context, id int64, success bool, status string) error
	Delete(ctx context.Context, id int64) error
}
