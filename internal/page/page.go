package page

import (
	"context"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/rs/zerolog"
)

type (
	ctxKeyType  int
	limitOffset int
)

const (
	offsetCtxKey ctxKeyType = iota
	limitCtxKey
)

const (
	defaultValueLimit  = 100
	defaultValueOffset = 0
)

// WithOffset returns context copy with offset value.
func WithOffset(ctx context.Context, offsetStr string) context.Context {
	logger := zerolog.Ctx(ctx)
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		logger.Err(err).Msg("Offset is missing or invalid, setting offset value to 0")
		offset = defaultValueOffset
	}
	return context.WithValue(ctx, offsetCtxKey, offset)
}

// WithLimit returns context copy with limit value.
func WithLimit(ctx context.Context, limitStr string) context.Context {
	logger := zerolog.Ctx(ctx)
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		logger.Err(err).Msg("Limit is missing or invalid, setting limit value to 100")
		limit = defaultValueLimit
	}
	return context.WithValue(ctx, limitCtxKey, limit)
}

func Limit(ctx context.Context) limitOffset {
	if lim, ok := ctx.Value(limitCtxKey).(int); ok {
		return limitOffset(lim)
	}
	return limitOffset(defaultValueLimit)
}

func Offset(ctx context.Context) limitOffset {
	if off, ok := ctx.Value(offsetCtxKey).(int); ok {
		return limitOffset(off)
	}
	return limitOffset(defaultValueOffset)
}

func (o limitOffset) IntPtr() *int {
	return ptr.To(int(o))
}

func (o limitOffset) Int() int {
	return int(o)
}

func (o limitOffset) Int64() int64 {
	return int64(o)
}

func (o limitOffset) String() string {
	return strconv.Itoa(int(o))
}
