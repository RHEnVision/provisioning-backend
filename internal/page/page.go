package page

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/math"

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
	tokenCtxKey
)

const (
	defaultValueLimit  = 100
	defaultValueOffset = 0
)

type Links struct {
	Previous string `json:"previous"`
	Next     string `json:"next"`
}
type Metadata struct {
	Total int   `json:"total,omitempty"`
	Links Links `json:"links,omitempty"`
}

// WithOffset returns context copy with offset value.
func WithOffset(ctx context.Context, offsetStr string) context.Context {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = defaultValueOffset
	}
	return context.WithValue(ctx, offsetCtxKey, offset)
}

// WithLimit returns context copy with limit value.
func WithLimit(ctx context.Context, limitStr string) context.Context {
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = defaultValueLimit
	}
	return context.WithValue(ctx, limitCtxKey, limit)
}

// WithToken returns context copy with token value.
func WithToken(ctx context.Context, token string) context.Context {
	logger := zerolog.Ctx(ctx)
	if token == "" {
		logger.Info().Msg("token is empty, listing first page")
	}
	return context.WithValue(ctx, tokenCtxKey, token)
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

func Token(ctx context.Context) string {
	if token, ok := ctx.Value(tokenCtxKey).(string); ok {
		return token
	}
	return ""
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

func (o limitOffset) Int32() int32 {
	return int32(o)
}

func (o limitOffset) String() string {
	return strconv.Itoa(int(o))
}

func NewOffsetMetadata(ctx context.Context, r *http.Request, total int) *Metadata {
	limit := Limit(ctx).Int()
	offset := Offset(ctx).Int()
	var prev, next string

	// First page: offset is 0
	if offset == 0 {
		prev = ""
	} else {
		prevOffset := math.Max(0, offset-limit)
		q := url.Values{}
		q.Add("limit", strconv.Itoa(limit))
		q.Add("offset", strconv.Itoa(prevOffset))
		prev = fmt.Sprintf("%v?%v", r.URL.Path, q.Encode())
	}

	// Last page: offset + Limit >= Total
	if offset+limit >= total {
		next = ""
	} else {
		nextOffset := offset + limit
		q := url.Values{}
		q.Add("limit", strconv.Itoa(limit))
		q.Add("offset", strconv.Itoa(nextOffset))
		next = fmt.Sprintf("%v?%v", r.URL.Path, q.Encode())
	}

	return &Metadata{
		Total: total,
		Links: Links{
			Next:     next,
			Previous: prev,
		},
	}
}

func NewTokenMetadata(ctx context.Context, r *http.Request, nextToken string) *Metadata {
	limit := Limit(ctx).Int()
	var next string

	if nextToken == "" {
		next = ""
	} else {
		q := url.Values{}
		q.Add("limit", strconv.Itoa(limit))
		q.Add("token", nextToken)
		next = fmt.Sprintf("%v?%v", r.URL.Path, q.Encode())
	}

	return &Metadata{
		Links: Links{
			Next: next,
		},
	}
}
