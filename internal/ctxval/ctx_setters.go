package ctxval

import (
	"context"

	"github.com/rs/zerolog"
)

func WithRequestID(parent context.Context, requestId string) context.Context {
	return context.WithValue(parent, requestIdCtxKey, requestId)
}

func WithRequestNum(parent context.Context, requestNum uint64) context.Context {
	return context.WithValue(parent, requestNumCtxKey, requestNum)
}

func WithLogger(parent context.Context, logger *zerolog.Logger) context.Context {
	return context.WithValue(parent, loggerCtxKey, logger)
}
