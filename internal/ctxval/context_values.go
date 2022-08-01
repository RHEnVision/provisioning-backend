package ctxval

import (
	"context"

	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type commonKeyId int

const (
	loggerCtxKey     commonKeyId = iota
	requestIdCtxKey  commonKeyId = iota
	requestNumCtxKey commonKeyId = iota
	accountIdCtxKey  commonKeyId = iota
)

// Identity returns identity header struct or nil when not set.
func Identity(ctx context.Context) identity.XRHID {
	return identity.Get(ctx)
}

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set.
func Logger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(loggerCtxKey) == nil {
		return &log.Logger
	}
	return ctx.Value(loggerCtxKey).(*zerolog.Logger)
}

func WithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// RequestId returns request id or an empty string when not set.
func RequestId(ctx context.Context) string {
	if ctx.Value(requestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(requestIdCtxKey).(string)
}

func WithRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, id)
}

// RequestNumber returns request counter.
func RequestNumber(ctx context.Context) uint64 {
	return ctx.Value(requestNumCtxKey).(uint64)
}

func WithRequestNumber(ctx context.Context, num uint64) context.Context {
	return context.WithValue(ctx, requestNumCtxKey, num)
}

// Account returns current account model or nil when not set.
func AccountId(ctx context.Context) int64 {
	return ctx.Value(accountIdCtxKey).(int64)
}

func WithAccountId(ctx context.Context, accountId int64) context.Context {
	return context.WithValue(ctx, accountIdCtxKey, accountId)
}
