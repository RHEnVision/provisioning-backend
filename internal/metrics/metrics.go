package metrics

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/prometheus/client_golang/prometheus"
)

var TotalSentAvailabilityCheckReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        "provisioning_source_availability_check_request_total",
		Help:        "availability check requests count partitioned by type (aws/gcp/azure), source status, component and error",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "statuser"},
	},
	[]string{"type", "status", "error"},
)

var TotalReceivedAvailabilityCheckReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        "provisioning_received_source_availability_check_request_total",
		Help:        "availability check requests count received from sources partitioned by type component and error",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "api"},
	},
	[]string{"error"},
)

var JobQueueSize = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
	Name:        "provisioning_job_queue_size",
	Help:        "background job queue size (pending tasks total)",
	ConstLabels: prometheus.Labels{"service": "provisioning"},
}, func() float64 {
	return float64(jq.Stats(context.Background()).EnqueuedJobs)
})

var AvailabilityCheckReqsDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:        "provisioning_source_availability_check_request_duration_ms",
		Help:        "availability check request duration partitioned by type and error",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "statuser"},
	},
	[]string{"type", "error"},
)

func ObserveAvailablilityCheckReqsDuration(provider models.ProviderType, ObservedFunc func() error) {
	errString := "false"
	start := time.Now()
	defer func() {
		AvailabilityCheckReqsDuration.WithLabelValues(provider.String(), errString).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}()

	err := ObservedFunc()
	if err != nil {
		errString = "true"
	}
}

func IncTotalSentAvailabilityCheckReqs(provider models.ProviderType, StatusType kafka.StatusType, err error) {
	errString := "false"
	if err != nil {
		errString = "true"
	}
	TotalSentAvailabilityCheckReqs.WithLabelValues(provider.String(), string(StatusType), errString).Inc()
}

func IncTotalReceivedAvailabilityCheckReqs(err error) {
	errString := "false"
	if err != nil {
		errString = "true"
	}
	TotalReceivedAvailabilityCheckReqs.WithLabelValues(errString).Inc()
}

func RegisterStatuserMetrics() {
	prometheus.MustRegister(TotalSentAvailabilityCheckReqs, AvailabilityCheckReqsDuration)
}

func RegisterApiMetrics() {
	prometheus.MustRegister(JobQueueSize, TotalReceivedAvailabilityCheckReqs)
}
