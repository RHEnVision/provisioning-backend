package metrics

import (
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/models"
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

var CacheHits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:        "provisioning_cache_hits",
	Help:        "The total number of cache hits per type with result (hit, miss, err)",
	ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
}, []string{"type", "result"})

var JobQueueSize = prometheus.NewGauge(prometheus.GaugeOpts{
	Name:        "provisioning_job_queue_size",
	Help:        "background job queue size (total pending jobs)",
	ConstLabels: prometheus.Labels{"service": "provisioning", "component": "stats"},
})

var JobQueueInFlight = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name:        "provisioning_job_queue_inflight",
	Help:        "number of in-flight jobs (total jobs which are currently processing)",
	ConstLabels: prometheus.Labels{"service": "provisioning", "component": "stats"},
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
		Name:        "provisioning_background_job_duration",
		Help:        "task queue job duration (in seconds) by type",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "worker"},
		Buckets:     []float64{0.5, 1, 2, 3, 4, 5, 7, 10, 30, 60 * 2, 60 * 10, 60 * 30},
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

var DbStatsDuration = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:        "provisioning_db_stats_duration",
		Help:        "task queue job duration (in ms)",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "stats"},
		Buckets:     []float64{10, 50, 100, 250, 500, 1000, 10000},
	},
)

var Reservations24hCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:        "provisioning_reservations_24h_count",
		Help:        "calculated sum of reservations in last 24 hours per result and provider",
		ConstLabels: prometheus.Labels{"service": "provisioning", "component": "stats"},
	},
	[]string{"result", "provider"},
)

var Reservations28dCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name:        "provisioning_reservations_28d_count",
		Help:        "calculated sum of reservations in last 28 days per result and provider",
		ConstLabels: prometheus.Labels{"service": "provisioning", "component": "stats"},
	},
	[]string{"result", "provider"},
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
		BackgroundJobDuration.WithLabelValues(jobType).Observe(time.Since(start).Seconds())
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

func IncCacheHit(model, result string) {
	CacheHits.WithLabelValues(model, result).Inc()
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

func ObserveDbStatsDuration(observedFunc func()) {
	start := time.Now()
	defer func() {
		DbStatsDuration.Observe(time.Since(start).Seconds())
	}()

	observedFunc()
}

func SetReservations24hCount(result string, pt models.ProviderType, count int64) {
	Reservations24hCount.WithLabelValues(result, pt.String()).Set(float64(count))
}

func SetReservations28dCount(result string, pt models.ProviderType, count int64) {
	Reservations28dCount.WithLabelValues(result, pt.String()).Set(float64(count))
}
