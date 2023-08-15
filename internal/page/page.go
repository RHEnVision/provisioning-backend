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
	Total int `json:"total"`
}

type Info struct {
	Metadata Metadata `json:"metadata"`
	Links    Links    `json:"links"`
}

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

func APIInfoResponse(ctx context.Context, r *http.Request, total int) *Info {
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

	return &Info{
		Metadata: Metadata{
			Total: total,
		},
		Links: Links{
			Next:     next,
			Previous: prev,
		},
	}
}
