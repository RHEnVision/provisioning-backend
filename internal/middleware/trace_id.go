package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/random"
	"go.opentelemetry.io/otel/trace"
)

func TraceID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		traceId := trace.SpanFromContext(ctx).SpanContext().TraceID()
		if !traceId.IsValid() {
			// OpenTelemetry library does not provide a public interface to create new IDs
			traceId = random.TraceID()
		}

		ctx = ctxval.WithTraceId(ctx, traceId.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
