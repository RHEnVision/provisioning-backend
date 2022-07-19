package ctxval

import (
	"context"

	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetRequestNum returns request number or a 0 when not in the context
func GetRequestNum(ctx context.Context) uint64 {
	if ctx.Value(requestNumCtxKey) == nil {
		return 0
	}
	return ctx.Value(requestNumCtxKey).(uint64)
}

// GetRequestId returns request id or an empty string when not in the context
func GetRequestId(ctx context.Context) string {
	if ctx.Value(requestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(requestIdCtxKey).(string)
}

// GetLogger returns logger or the standard global logger when not in the context
// or when context is nil.
func GetLogger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(loggerCtxKey) == nil {
		return &log.Logger
	}
	logger := ctx.Value(loggerCtxKey).(*zerolog.Logger)
	return logger
}

func GetIdentity(ctx context.Context) identity.XRHID {
	return identity.Get(ctx)
}
