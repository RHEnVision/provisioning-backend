package ctxval

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetStringValue(ctx context.Context, key CommonKeyId) string {
	return ctx.Value(key).(string)
}

func GetUInt64Value(ctx context.Context, key CommonKeyId) uint64 {
	return ctx.Value(key).(uint64)
}

// GetLogger returns logger or the standard global logger when not in the context
// or when context is nil.
func GetLogger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(LoggerCtxKey) == nil {
		return &log.Logger
	}
	logger := ctx.Value(LoggerCtxKey).(zerolog.Logger)
	return &logger
}

// GetRequestId returns request id or an empty string when not in the context
func GetRequestId(ctx context.Context) string {
	if ctx.Value(RequestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(RequestIdCtxKey).(string)
}
