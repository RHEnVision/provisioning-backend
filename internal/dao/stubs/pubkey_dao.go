package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type pubkeyCtxKeyType string

var pubkeyCtxKey pubkeyCtxKeyType = "pubkey-dao"

type PubkeyDaoStub struct {
	store []*models.Pubkey
}

func init() {
	dao.GetPubkeyDao = getPubkeyDao
}

type ContextReadError struct{}

func (m *ContextReadError) Error() string {
	return "failed to find or convert dao stored in testing context"
}

func WithPubkeyDao(parent context.Context, init_store []*models.Pubkey) context.Context {
	ctx := context.WithValue(parent, pubkeyCtxKey, PubkeyDaoStub{init_store})
	return ctx
}

func getPubkeyDao(ctx context.Context) (dao.PubkeyDao, error) {
	var err error
	pkdao, ok := ctx.Value(pubkeyCtxKey).(dao.PubkeyDao)
	if !ok {
		err = &ContextReadError{}
	}

	return pkdao, err
}

func (PubkeyDaoStub) Create(ctx context.Context, pk *models.Pubkey) error { return nil }
func (PubkeyDaoStub) CreateWithResource(ctx context.Context, pk *models.Pubkey, pkr *models.PubkeyResource) error {
	return nil
}
func (PubkeyDaoStub) Update(ctx context.Context, pk *models.Pubkey) error { return nil }
func (mock PubkeyDaoStub) GetById(ctx context.Context, id uint64) (*models.Pubkey, error) {
	return mock.store[0], nil
}
func (mock PubkeyDaoStub) List(ctx context.Context, limit, offset uint64) ([]*models.Pubkey, error) {
	return mock.store, nil
}
func (PubkeyDaoStub) Delete(ctx context.Context, id uint64) error { return nil }
