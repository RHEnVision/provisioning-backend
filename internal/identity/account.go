package identity

import (
	"context"
	"errors"
)

type cxtKeyId int

const (
	accountIdCtxKey cxtKeyId = iota
)

var MissingAccountInContextError = errors.New("operation requires account_id in context")

// AccountId returns current account model or panics when not set
func AccountId(ctx context.Context) int64 {
	value := ctx.Value(accountIdCtxKey)
	if value == nil {
		panic(MissingAccountInContextError)
	}
	return value.(int64)
}

// AccountIdOrNil returns current account model or 0 when not set.
func AccountIdOrNil(ctx context.Context) int64 {
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
