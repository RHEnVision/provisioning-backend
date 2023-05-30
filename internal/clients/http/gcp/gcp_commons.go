package gcp

import (
	"context"

	"github.com/rs/zerolog"
)

func logger(ctx context.Context) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("client", "gcp").Logger()
}
