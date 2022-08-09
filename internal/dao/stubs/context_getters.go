package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type daoStubCtxKeyType int

const (
	accountCtxKey daoStubCtxKeyType = iota
	pubkeyCtxKey  daoStubCtxKeyType = iota
)

func ctxAccountId(ctx context.Context) int64 {
	return ctxval.AccountId(ctx)
}

func WithPubkeyDao(parent context.Context) context.Context {
	if parent.Value(pubkeyCtxKey) != nil {
		panic(ContextSecondInitializationError)
	}

	ctx := context.WithValue(parent, pubkeyCtxKey, &pubkeyDaoStub{lastId: 0, store: []*models.Pubkey{}})
	return ctx
}

func getPubkeyDaoStub(ctx context.Context) (*pubkeyDaoStub, error) {
	var ok bool
	var pkdao *pubkeyDaoStub
	if pkdao, ok = ctx.Value(pubkeyCtxKey).(*pubkeyDaoStub); !ok {
		return nil, ContextReadError
	}
	return pkdao, nil
}

func WithAccountDaoOne(parent context.Context) context.Context {
	if parent.Value(accountCtxKey) != nil {
		panic(ContextSecondInitializationError)
	}

	ctx := context.WithValue(parent, accountCtxKey, buildAccountDaoWithOneAccount())
	return ctx
}

func WithAccountDaoNull(parent context.Context) context.Context {
	if parent.Value(accountCtxKey) != nil {
		panic(ContextSecondInitializationError)
	}

	ctx := context.WithValue(parent, accountCtxKey, buildAccountDaoWithNullValue())
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
