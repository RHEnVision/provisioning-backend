package ctxval

import (
	"context"

	ucontext "github.com/Unleash/unleash-client-go/v3/context"
)

// UnleashContext returns unleash context or an empty context when not set.
func UnleashContext(ctx context.Context) ucontext.Context {
	if ctx.Value(unleashContextCtxKey) == nil {
		return ucontext.Context{}
	}
	return ctx.Value(unleashContextCtxKey).(ucontext.Context)
}

// WithUnleashContext returns context copy with unleash context as a value.
func WithUnleashContext(ctx context.Context, uctx ucontext.Context) context.Context {
	return context.WithValue(ctx, unleashContextCtxKey, uctx)
}
