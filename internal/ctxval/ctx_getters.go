package ctxval

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetStringValue(ctx context.Context, key CommonKeyId) string {
	return ctx.Value(key).(string)
}

func GetUInt64Value(ctx context.Context, key CommonKeyId) uint64 {
	return ctx.Value(key).(uint64)
}

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set.
func Logger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(LoggerCtxKey) == nil {
		return &log.Logger
	}
	logger := ctx.Value(LoggerCtxKey).(zerolog.Logger)
	return &logger
}

// RequestId returns request id or an empty string when not set.
func RequestId(ctx context.Context) string {
	if ctx.Value(RequestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(RequestIdCtxKey).(string)
}

// Identity returns identity header struct or nil when not set.
func Identity(ctx context.Context) identity.XRHID {
	return identity.Get(ctx)
}

// Account returns current account model or nil when not set.
func Account(ctx context.Context) *models.Account {
	return ctx.Value(AccountCtxKey).(*models.Account)
}
