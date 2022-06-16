package dao

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

var GetAccountDao func(ctx context.Context) (AccountDao, error)

type AccountDao interface {
	GetById(ctx context.Context, id int64) (*models.Account, error)
	GetByAccountNumber(ctx context.Context, number string) (*models.Account, error)
	GetByOrgId(ctx context.Context, orgId string) (*models.Account, error)
	List(ctx context.Context, limit, offset int64) ([]*models.Account, error)
}

var GetPubkeyDao func(ctx context.Context) (PubkeyDao, error)

type PubkeyDao interface {
	Create(ctx context.Context, pk *models.Pubkey) error
	CreateWithResource(ctx context.Context, pk *models.Pubkey, pkr *models.PubkeyResource) error
	Update(ctx context.Context, pk *models.Pubkey) error
	GetById(ctx context.Context, id int64) (*models.Pubkey, error)
	List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error)
	Delete(ctx context.Context, id int64) error
}

var GetPubkeyResourceDao func(ctx context.Context) (PubkeyResourceDao, error)

type PubkeyResourceDao interface {
	ListByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error)
	Create(ctx context.Context, pkr *models.PubkeyResource) error
	Update(ctx context.Context, pkr *models.PubkeyResource) error
	Delete(ctx context.Context, id int64) error
}
