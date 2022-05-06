package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var Buckets []float64 = []float64{.05, .1, .25, .5, 1}

var HTTPRequestsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "How many HTTP requests processed, partitioned by status code, http method and path.",
	},
	[]string{"code", "method", "path"},
)

var HTTPResponseTimeSec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_response_time_seconds",
	Help:    "Duration of HTTP requests, partitioned by path.",
	Buckets: Buckets,
}, []string{"path"})

func init() {
	prometheus.MustRegister(HTTPRequestsTotal)
	prometheus.MustRegister(HTTPResponseTimeSec)
}
