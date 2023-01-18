package ctxval

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set. Never returns nil.
func Logger(ctx context.Context) *zerolog.Logger {
	value := ctx.Value(loggerCtxKey)
	if ctx == nil || value == nil {
		return &log.Logger
	}
	return value.(*zerolog.Logger)
}

// WithLogger returns context copy with logger.
func WithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
