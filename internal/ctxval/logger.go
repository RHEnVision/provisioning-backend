package ctxval

import (
	"context"

	"github.com/rs/zerolog"
)

// Logger returns the main logger with context fields or the standard global logger
// when the main logger was not set. When nil is passed, returns noop logger.
// Never returns nil.
//
// Deprecated: Use zerolog.Ctx or log.Ctx alias instead.
func Logger(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// WithLogger returns context copy with logger.
//
// Deprecated: Use log.Logger.WithContext(ctx) instead.
func WithLogger(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}
