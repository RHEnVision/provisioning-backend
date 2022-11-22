package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type pubkeyDaoStub struct {
	lastId int64
	store  []*models.Pubkey
}

func init() {
	dao.GetPubkeyDao = getPubkeyDao
}

func PubkeyStubCount(ctx context.Context) int {
	pkdao := getPubkeyDaoStub(ctx)
	return len(pkdao.store)
}

func getPubkeyDao(ctx context.Context) dao.PubkeyDao {
	return getPubkeyDaoStub(ctx)
}

func (stub *pubkeyDaoStub) Create(ctx context.Context, pubkey *models.Pubkey) error {
	if pubkey.AccountID == 0 {
		pubkey.AccountID = ctxAccountId(ctx)
	}
	if pubkey.AccountID != ctxAccountId(ctx) {
		return dao.ErrWrongAccount
	}
	if err := models.Validate(ctx, pubkey); err != nil {
		return dao.ErrValidation
	}
	if err := models.Transform(ctx, pubkey); err != nil {
		return dao.ErrTransformation
	}

	pubkey.ID = stub.lastId + 1
	stub.store = append(stub.store, pubkey)
	stub.lastId++
	return nil
}

func (stub *pubkeyDaoStub) Update(ctx context.Context, pubkey *models.Pubkey) error {
	if pubkey.AccountID == 0 {
		pubkey.AccountID = ctxAccountId(ctx)
	}
	if pubkey.AccountID != ctxAccountId(ctx) {
		return dao.ErrWrongAccount
	}

	for idx, p := range stub.store {
		if p.ID == pubkey.ID {
			stub.store[idx] = pubkey
			return nil
		}
	}
	return dao.ErrNoRows
}

func (stub *pubkeyDaoStub) GetById(ctx context.Context, id int64) (*models.Pubkey, error) {
	for _, pk := range stub.store {
		if pk.AccountID == ctxAccountId(ctx) && pk.ID == id {
			return pk, nil
		}
	}
	return nil, dao.ErrNoRows
}

func (stub *pubkeyDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error) {
	var filtered []*models.Pubkey
	for _, pk := range stub.store {
		if pk.AccountID == ctxAccountId(ctx) {
			filtered = append(filtered, pk)
		}
	}
	return filtered, nil
}

func (stub *pubkeyDaoStub) Delete(ctx context.Context, id int64) error {
	for idx, p := range stub.store {
		if p.AccountID == ctxAccountId(ctx) && p.ID == id {
			stub.store = append(stub.store[:idx], stub.store[idx+1:]...)
			return nil
		}
	}
	return nil
}

func (stub *pubkeyDaoStub) UnscopedGetResourceBySourceAndRegion(ctx context.Context, pubkeyId int64, sourceId string, region string) (*models.PubkeyResource, error) {
	return nil, nil
}

func (stub *pubkeyDaoStub) UnscopedCreateResource(ctx context.Context, pkr *models.PubkeyResource) error {
	return nil
}

func (stub *pubkeyDaoStub) UnscopedDeleteResource(ctx context.Context, id int64) error {
	return nil
}

func (stub *pubkeyDaoStub) UnscopedListResourcesByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error) {
	return nil, nil
}
