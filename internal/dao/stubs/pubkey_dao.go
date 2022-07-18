package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type PubkeyDaoStub struct {
	store []*models.Pubkey
}

func init() {
	dao.GetPubkeyDao = getPubkeyDao
}

func PubkeyStubCount(ctx context.Context) (int, error) {
	pkdao, err := getPubkeyDaoStub(ctx)
	if err != nil {
		return 0, err
	}
	return len(pkdao.store), nil
}

func getPubkeyDao(ctx context.Context) (dao.PubkeyDao, error) {
	return getPubkeyDaoStub(ctx)
}

func (mock *PubkeyDaoStub) Create(ctx context.Context, pk *models.Pubkey) error {
	mock.store = append(mock.store, pk)
	return nil
}
func (*PubkeyDaoStub) CreateWithResource(ctx context.Context, pk *models.Pubkey, pkr *models.PubkeyResource) error {
	return nil
}
func (*PubkeyDaoStub) Update(ctx context.Context, pk *models.Pubkey) error { return nil }
func (mock *PubkeyDaoStub) GetById(ctx context.Context, id int64) (*models.Pubkey, error) {
	return mock.store[0], nil
}
func (mock *PubkeyDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error) {
	return mock.store, nil
}
func (*PubkeyDaoStub) Delete(ctx context.Context, id int64) error { return nil }
