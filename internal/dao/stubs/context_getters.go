package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type daoStubCtxKeyType int

const (
	accountCtxKey daoStubCtxKeyType = iota
	pubkeyCtxKey  daoStubCtxKeyType = iota
)

func WithPubkeyDao(parent context.Context) context.Context {
	ctx := context.WithValue(parent, pubkeyCtxKey, &PubkeyDaoStub{lastId: 0, store: []*models.Pubkey{}})
	return ctx
}

func getPubkeyDaoStub(ctx context.Context) (*PubkeyDaoStub, error) {
	var ok bool
	var pkdao *PubkeyDaoStub
	if pkdao, ok = ctx.Value(pubkeyCtxKey).(*PubkeyDaoStub); !ok {
		return nil, ContextReadError
	}
	return pkdao, nil
}

func WithAccountDaoOne(parent context.Context) context.Context {
	ctx := context.WithValue(parent, accountCtxKey, buildAccountDaoWithOneAccount())
	return ctx
}

func getAccountDaoStub(ctx context.Context) (*accountDaoStub, error) {
	var ok bool
	var accdao *accountDaoStub
	if accdao, ok = ctx.Value(accountCtxKey).(*accountDaoStub); !ok {
		return nil, ContextReadError
	}
	return accdao, nil
}
