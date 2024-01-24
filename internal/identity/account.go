package identity

import (
	"context"
	"errors"
)

type cxtKeyId int

const (
	accountIdCtxKey cxtKeyId = iota
)

var ErrMissingAccountInContext = errors.New("operation requires account_id in context")

// AccountId returns current account model or panics when not set
func AccountId(ctx context.Context) int64 {
	value := ctx.Value(accountIdCtxKey)
	if value == nil {
		panic(ErrMissingAccountInContext)
	}
	return value.(int64)
}

// AccountIdOrZero returns current account model or 0 when not set.
func AccountIdOrZero(ctx context.Context) int64 {
	value := ctx.Value(accountIdCtxKey)
	if value == nil {
		return 0
	}
	return value.(int64)
}

// WithAccountId returns context copy with account id value.
func WithAccountId(ctx context.Context, accountId int64) context.Context {
	return context.WithValue(ctx, accountIdCtxKey, accountId)
}
