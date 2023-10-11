package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const TraceName = telemetry.TracePrefix + "internal/middleware"

// Telemetry middleware starts a new telemetry span for this request,
// it tries to find the parent trace in the request,
// if none is found, it starts new root span.
func Telemetry(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var span trace.Span
		ctx := r.Context()

		ctx, span = otel.Tracer(TraceName).Start(ctx, r.URL.Path)

		// Store TraceID in response headers for easier debugging
		w.Header().Set("X-Trace-Id", span.SpanContext().TraceID().String())

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
