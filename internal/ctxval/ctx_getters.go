package ctxval

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set.
func Logger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(loggerCtxKey) == nil {
		return &log.Logger
	}
	// TODO: we should store pointer to logger instead: https://issues.redhat.com/browse/HMSPROV-118
	logger := ctx.Value(loggerCtxKey).(zerolog.Logger)
	return &logger
}

func SetLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

// RequestId returns request id or an empty string when not set.
func RequestId(ctx context.Context) string {
	if ctx.Value(requestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(requestIdCtxKey).(string)
}

func SetRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, id)
}

// RequestNumber returns request counter.
func RequestNumber(ctx context.Context) uint64 {
	return ctx.Value(requestNumCtxKey).(uint64)
}

func SetRequestNumber(ctx context.Context, num uint64) context.Context {
	return context.WithValue(ctx, requestNumCtxKey, num)
}

// Identity returns identity header struct or nil when not set.
func Identity(ctx context.Context) identity.XRHID {
	return identity.Get(ctx)
}

// Account returns current account model or nil when not set.
func Account(ctx context.Context) *models.Account {
	return ctx.Value(accountCtxKey).(*models.Account)
}

func SetAccount(ctx context.Context, account *models.Account) context.Context {
	return context.WithValue(ctx, accountCtxKey, account)
}
