package metrics

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue/dejq"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/prometheus/client_golang/prometheus"
)

var TotalAvailabilityCheckReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        "provisioning_source_availability_check_request_total",
		Help:        "availability Check requests count partitioned by type (aws/gcp/azure), source status, component and error",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "statuser"},
	},
	[]string{"type", "status", "error"},
)

var JobQueueSize = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
	Name:        "provisioning_job_queue_size",
	Help:        "background job queue size (pending tasks total)",
	ConstLabels: prometheus.Labels{"service": "provisioning"},
}, func() float64 {
	return float64(dejq.Stats(context.Background()).EnqueuedJobs)
})

var AvailabilityCheckReqsDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:        "provisioning_source_availability_check_request_duration_ms",
		Help:        "Availability check Request duration partitioned by type and error",
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

func IncTotalAvailabilityCheckReqs(provider models.ProviderType, StatusType kafka.StatusType, err error) {
	errString := "false"
	if err != nil {
		errString = "true"
	}
	TotalAvailabilityCheckReqs.WithLabelValues(provider.String(), string(StatusType), errString).Inc()
}

func RegisterStatuserMetrics() {
	prometheus.MustRegister(TotalAvailabilityCheckReqs, AvailabilityCheckReqsDuration)
}

func RegisterApiMetrics() {
	prometheus.MustRegister(JobQueueSize)
}
