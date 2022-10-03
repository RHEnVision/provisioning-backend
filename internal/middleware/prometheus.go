package middleware

// Chi routes aware middleware, taken from https://github.com/766b/chi-prometheus
// License: Apache 2.0

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var buckets = []float64{100, 200, 500, 5000}

const (
	metricNameHttpRequestTotal    = "http_request_total"
	metricNameHttpRequestDuration = "http_request_duration_ms"
)

// Middleware is a handler that exposes prometheus metrics for the number of requests,
// the latency and the response size, partitioned by status code, method and HTTP path.
type Middleware struct {
	reqs    *prometheus.CounterVec
	latency *prometheus.HistogramVec
}

// NewPatternMiddleware returns a new prometheus Middleware handler that groups requests by the chi routing pattern.
// EX: /users/{firstName} instead of /users/bob
func NewPatternMiddleware(name string) func(next http.Handler) http.Handler {
	var m Middleware
	m.reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        metricNameHttpRequestTotal,
			Help:        "HTTP requests count partitioned by status code, method and HTTP path (chi route)",
			ConstLabels: prometheus.Labels{"service": name},
		},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.reqs)

	m.latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        metricNameHttpRequestDuration,
		Help:        "Request duration partitioned by status code, method and HTTP path (chi route)",
		ConstLabels: prometheus.Labels{"service": name},
		Buckets:     buckets,
	},
		[]string{"code", "method", "path"},
	)
	prometheus.MustRegister(m.latency)
	return m.patternHandler
}

func (c Middleware) patternHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		rctx := chi.RouteContext(r.Context())
		routePattern := strings.Join(rctx.RoutePatterns, "")
		routePattern = strings.Replace(routePattern, "/*/", "/", -1)

		c.reqs.WithLabelValues(http.StatusText(ww.Status()), r.Method, routePattern).Inc()
		c.latency.WithLabelValues(http.StatusText(ww.Status()), r.Method, routePattern).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}
	return http.HandlerFunc(fn)
}
