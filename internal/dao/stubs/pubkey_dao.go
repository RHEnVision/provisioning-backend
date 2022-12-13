package stubs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type pubkeyDaoStub struct {
	lastId        int64
	store         []*models.Pubkey
	resourceStore []*models.PubkeyResource
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
	var resultPk models.Pubkey
	for _, pk := range stub.store {
		if pk.AccountID == ctxAccountId(ctx) && pk.ID == id {
			resultPk = *pk // shallow copy of the struct
			return &resultPk, nil
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
	for _, pkr := range stub.resourceStore {
		if pkr.PubkeyID == pubkeyId && pkr.SourceID == sourceId && pkr.Region == region {
			return pkr, nil
		}
	}
	return nil, dao.ErrNoRows
}

func (stub *pubkeyDaoStub) UnscopedCreateResource(ctx context.Context, pkr *models.PubkeyResource) error {
	pkr.ID = int64(len(stub.resourceStore)) + 1
	stub.resourceStore = append(stub.resourceStore, pkr)
	return nil
}

func (stub *pubkeyDaoStub) UnscopedDeleteResource(ctx context.Context, id int64) error {
	return nil
}

func (stub *pubkeyDaoStub) UnscopedListResourcesByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error) {
	var result []*models.PubkeyResource
	for _, pkr := range stub.resourceStore {
		if pkr.PubkeyID == pkId {
			shallowCopy := *pkr
			result = append(result, &shallowCopy)
		}
	}
	return result, nil
}

func (stub *pubkeyDaoStub) UnscopedUpdateHandle(ctx context.Context, id int64, handle string) error {
	for _, pkr := range stub.resourceStore {
		if pkr.ID == id {
			pkr.Handle = handle
			return nil
		}
	}
	return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
}
