package ctxval

import "context"

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
