package logging

import (
	"context"
)

type commonKeyId int

const (
	requestIdCtxKey     commonKeyId = iota
	edgeRequestIdCtxKey commonKeyId = iota
	correlationCtxKey   commonKeyId = iota
)

// CorrelationId returns UI correlation id or an empty string when not set.
func CorrelationId(ctx context.Context) string {
	value := ctx.Value(correlationCtxKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

// WithCorrelationId returns context copy with correlation id value.
func WithCorrelationId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationCtxKey, id)
}

// EdgeRequestId returns edge API (3Scale) request id or an empty string when not set.
func EdgeRequestId(ctx context.Context) string {
	value := ctx.Value(edgeRequestIdCtxKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

// WithEdgeRequestId returns context copy with trace id value.
func WithEdgeRequestId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, edgeRequestIdCtxKey, id)
}

// TraceId returns request id or an empty string when not set.
func TraceId(ctx context.Context) string {
	value := ctx.Value(requestIdCtxKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

// WithTraceId returns context copy with trace id value.
func WithTraceId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, id)
}
