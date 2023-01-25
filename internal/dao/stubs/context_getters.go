package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type daoStubCtxKeyType int

const (
	accountCtxKey     daoStubCtxKeyType = iota
	pubkeyCtxKey      daoStubCtxKeyType = iota
	reservationCtxKey daoStubCtxKeyType = iota
)

func ctxAccountId(ctx context.Context) int64 {
	return ctxval.AccountId(ctx)
}

func WithPubkeyDao(parent context.Context) context.Context {
	if parent.Value(pubkeyCtxKey) != nil {
		panic(dao.ErrStubContextAlreadySet)
	}

	ctx := context.WithValue(parent, pubkeyCtxKey, &pubkeyDaoStub{lastId: 0, store: []*models.Pubkey{}})
	return ctx
}

func getPubkeyDaoStub(ctx context.Context) *pubkeyDaoStub {
	var ok bool
	var pkdao *pubkeyDaoStub
	if pkdao, ok = ctx.Value(pubkeyCtxKey).(*pubkeyDaoStub); !ok {
		panic(dao.ErrStubMissingContext)
	}
	return pkdao
}

func WithReservationDao(parent context.Context) context.Context {
	if parent.Value(reservationCtxKey) != nil {
		panic(dao.ErrStubContextAlreadySet)
	}

	ctx := context.WithValue(parent, reservationCtxKey, &reservationDaoStub{})
	return ctx
}

func getReservationDaoStub(ctx context.Context) *reservationDaoStub {
	var ok bool
	var resDao *reservationDaoStub
	if resDao, ok = ctx.Value(reservationCtxKey).(*reservationDaoStub); !ok {
		panic(dao.ErrStubMissingContext)
	}
	return resDao
}

func WithAccountDaoOne(parent context.Context) context.Context {
	if parent.Value(accountCtxKey) != nil {
		panic(dao.ErrStubContextAlreadySet)
	}

	ctx := context.WithValue(parent, accountCtxKey, buildAccountDaoWithOneAccount())
	return ctx
}

func WithAccountDaoNull(parent context.Context) context.Context {
	if parent.Value(accountCtxKey) != nil {
		panic(dao.ErrStubContextAlreadySet)
	}

	ctx := context.WithValue(parent, accountCtxKey, buildAccountDaoWithNullValue())
	return ctx
}

func getAccountDaoStub(ctx context.Context) *accountDaoStub {
	var ok bool
	var accdao *accountDaoStub
	if accdao, ok = ctx.Value(accountCtxKey).(*accountDaoStub); !ok {
		panic(dao.ErrStubMissingContext)
	}
	return accdao
}
