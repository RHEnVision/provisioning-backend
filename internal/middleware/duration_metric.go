package middleware

import (
	"net/http"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}
		timer := prometheus.NewTimer(metrics.HTTPResponseTimeSec.WithLabelValues(r.URL.Path))
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode
		metrics.HTTPRequestsTotal.WithLabelValues(strconv.Itoa(statusCode), r.Method, r.URL.Path).Inc()
		timer.ObserveDuration()
	})
}
