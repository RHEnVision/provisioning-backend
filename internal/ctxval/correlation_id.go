package ctxval

import "context"

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
