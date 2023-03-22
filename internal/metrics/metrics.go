package metrics

import (
	"time"

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

var TotalInvalidAvailabilityCheckReqs = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name:        "provisioning_invalid_source_availability_check_request_total",
		Help:        "invalid availability check requests count partitioned by component",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "statuser"},
	},
)

var JobQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
	Name:        "provisioning_job_queue_size",
	Help:        "background job queue size (total pending jobs)",
	ConstLabels: prometheus.Labels{"service": "provisioning"},
})

var JobQueueInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:        "provisioning_job_queue_inflight",
	Help:        "number of in-flight jobs (total jobs which are currently processing)",
	ConstLabels: prometheus.Labels{"service": "provisioning"},
}, []string{"worker"})

var AvailabilityCheckReqsDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:        "provisioning_source_availability_check_request_duration_ms",
		Help:        "availability check request duration partitioned by type and error",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "statuser"},
		Buckets:     []float64{50, 100, 250, 500, 1000, 2500, 6000},
	},
	[]string{"type", "error"},
)

var BackgroundJobDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:        "provisioning_background_job_duration_ms",
		Help:        "task queue job duration (ms) by type",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "worker"},
		Buckets:     []float64{50, 100, 250, 500, 1000, 2500, 6000},
	},
	[]string{"type"},
)

var ReservationCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        "provisioning_reservation_count",
		Help:        "reservation count by result (success/failure) by type (aws/gcp/azure)",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "worker"},
	},
	[]string{"type", "result"},
)

func ObserveAvailabilityCheckReqsDuration(provider string, observedFunc func() error) {
	errString := "false"
	start := time.Now()
	defer func() {
		AvailabilityCheckReqsDuration.WithLabelValues(provider, errString).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}()

	err := observedFunc()
	if err != nil {
		errString = "true"
	}
}

func ObserveBackgroundJobDuration(jobType string, observedFunc func()) {
	start := time.Now()
	defer func() {
		BackgroundJobDuration.WithLabelValues(jobType).Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}()

	observedFunc()
}

func IncTotalSentAvailabilityCheckReqs(provider string, statusType string, err error) {
	errString := "false"
	if err != nil {
		errString = "true"
	}
	TotalSentAvailabilityCheckReqs.WithLabelValues(provider, string(statusType), errString).Inc()
}

func IncTotalInvalidAvailabilityCheckReqs() {
	TotalInvalidAvailabilityCheckReqs.Inc()
}

func SetJobQueueSize(size uint64) {
	JobQueueSize.Set(float64(size))
}

func SetJobQueueInFlight(workerName string, inflight int64) {
	JobQueueInFlight.WithLabelValues(workerName).Set(float64(inflight))
}

func IncReservationCount(rtype, result string) {
	ReservationCount.WithLabelValues(rtype, result).Inc()
}
