package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type PubkeyDaoStub struct {
	lastId int64
	store  []*models.Pubkey
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

func (stub *PubkeyDaoStub) Create(ctx context.Context, pk *models.Pubkey) error {
	pk.ID = stub.lastId + 1
	stub.store = append(stub.store, pk)
	stub.lastId++
	return nil
}

func (stub *PubkeyDaoStub) Update(ctx context.Context, pk *models.Pubkey) error {
	for idx, p := range stub.store {
		if p.ID == pk.ID {
			stub.store[idx] = pk
			return nil
		}
	}
	return NewRecordNotFoundError(ctx, "Pubkey")
}

func (stub *PubkeyDaoStub) GetById(ctx context.Context, id int64) (*models.Pubkey, error) {
	for _, acc := range stub.store {
		if acc.ID == id {
			return acc, nil
		}
	}
	return nil, NewRecordNotFoundError(ctx, "Pubkey")
}

func (stub *PubkeyDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error) {
	return stub.store, nil
}

func (stub *PubkeyDaoStub) Delete(ctx context.Context, id int64) error {
	for idx, p := range stub.store {
		if p.ID == id {
			stub.store = append(stub.store[:idx], stub.store[idx+1:]...)
			return nil
		}
	}
	return nil
}
