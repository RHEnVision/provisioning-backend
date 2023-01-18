package ctxval

import "context"

// TraceId returns request id or an empty string when not set.
func TraceId(ctx context.Context) string {
	if ctx.Value(requestIdCtxKey) == nil {
		return ""
	}
	return ctx.Value(requestIdCtxKey).(string)
}

// WithTraceId returns context copy with trace id value.
func WithTraceId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, id)
}
