package ctxval

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set. Never returns nil.
func Logger(ctx context.Context) *zerolog.Logger {
	if ctx == nil || ctx.Value(loggerCtxKey) == nil {
		return &log.Logger
	}
	return ctx.Value(loggerCtxKey).(*zerolog.Logger)
}

// WithLogger returns context copy with logger.
func WithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
