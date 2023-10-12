package logging

import (
	"context"
)

type commonKeyId int

const (
	requestIdCtxKey     commonKeyId = iota
	edgeRequestIdCtxKey commonKeyId = iota
	correlationCtxKey   commonKeyId = iota
	jobIdCtxKey         commonKeyId = iota
	reservationIdCtxKey commonKeyId = iota
	jobTypeCtxKey       commonKeyId = iota
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

// JobId returns request id or an empty string when not set.
func JobId(ctx context.Context) string {
	value := ctx.Value(jobIdCtxKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

// WithJobId returns context copy with trace id value.
func WithJobId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, jobIdCtxKey, id)
}

// ReservationId returns request id or an empty string when not set.
func ReservationId(ctx context.Context) int64 {
	value := ctx.Value(reservationIdCtxKey)
	if value == nil {
		return 0
	}
	return value.(int64)
}

// WithReservationId returns context copy with trace id value.
func WithReservationId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, reservationIdCtxKey, id)
}

// JobType returns relevant context data or empty string when not set.
func JobType(ctx context.Context) string {
	value := ctx.Value(jobTypeCtxKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

// WithJobType returns context copy with relevant value.
func WithJobType(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, jobTypeCtxKey, id)
}
