package dao

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

var GetAccountDao func(ctx context.Context) (AccountDao, error)

type AccountDao interface {
	GetById(ctx context.Context, id uint64) (*models.Account, error)
	GetByAccountNumber(ctx context.Context, number string) (*models.Account, error)
	GetByOrgId(ctx context.Context, orgId string) (*models.Account, error)
	List(ctx context.Context, limit, offset uint64) ([]*models.Account, error)
}
