package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/smithy-go/ptr"
)

type accountDaoStub struct {
	store  []*models.Account
	lastId int64
}

func buildAccountDaoWithOneAccount() *accountDaoStub {
	return &accountDaoStub{
		lastId: 1,
		store: []*models.Account{{
			ID:            1,
			OrgID:         "1",
			AccountNumber: ptr.String("1"),
		}},
	}
}

func init() {
	dao.GetAccountDao = getAccountDao
}

func getAccountDao(ctx context.Context) (dao.AccountDao, error) {
	return getAccountDaoStub(ctx)
}

func (stub *accountDaoStub) Create(ctx context.Context, pk *models.Account) error {
	pk.ID = stub.lastId + 1
	stub.store = append(stub.store, pk)
	stub.lastId++
	return nil
}

func (stub *accountDaoStub) GetById(ctx context.Context, id int64) (*models.Account, error) {
	for _, acc := range stub.store {
		if acc.ID == id {
			return acc, nil
		}
	}
	return nil, NewRecordNotFoundError(ctx, "Account")
}

func (stub *accountDaoStub) GetOrCreateByIdentity(ctx context.Context, orgId string, accountNumber string) (*models.Account, error) {
	acc, err := stub.GetByOrgId(ctx, orgId)
	if err == nil {
		return acc, nil
	}
	acc, err = stub.GetByAccountNumber(ctx, accountNumber)
	if err == nil {
		return acc, nil
	}
	return nil, NewRecordNotFoundError(ctx, "Account")
}

func (stub *accountDaoStub) GetByOrgId(ctx context.Context, orgId string) (*models.Account, error) {
	for _, acc := range stub.store {
		if acc.OrgID == orgId {
			return acc, nil
		}
	}
	return nil, NewRecordNotFoundError(ctx, "Account")
}

func (stub *accountDaoStub) GetByAccountNumber(ctx context.Context, number string) (*models.Account, error) {
	for _, acc := range stub.store {
		if *acc.AccountNumber == number {
			return acc, nil
		}
	}
	return nil, NewRecordNotFoundError(ctx, "Account")
}

func (stub *accountDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Account, error) {
	return stub.store, nil
}
