package metrics

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue/dejq"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

var TotalAvailabilityCheckReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        "source_availability_check_request_total",
		Help:        "availability Check requests count partitioned by type (aws/gcp/azure), source status, component and error",
		ConstLabels: prometheus.Labels{"service": "provisioning"},
	},
	[]string{"type", "status", "component", "error"},
)

var JobQueueSize = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
	Name:        "job_queue_size",
	Help:        "background job queue size (pending tasks total)",
	ConstLabels: prometheus.Labels{"service": "provisioning"},
}, func() float64 {
	return float64(dejq.Stats(context.Background()).EnqueuedJobs)
})

func IncTotalAvailabilityCheckReqs(provider models.ProviderType, component string, StatusType kafka.StatusType, err error) {
	errString := "false"
	if err != nil {
		errString = "true"
	}
	TotalAvailabilityCheckReqs.WithLabelValues(provider.String(), string(StatusType), component, errString).Inc()
}

func RegisterStatuserMetrics() {
	prometheus.MustRegister(TotalAvailabilityCheckReqs)
}

func RegisterApiMetrics() {
	prometheus.MustRegister(JobQueueSize)
}
