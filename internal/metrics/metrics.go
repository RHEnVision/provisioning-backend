package metrics

import (
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricNameSourceAvailabilityCheckRequestTotal = "source_availability_check_request_total"
)

var TotalAvailabilityCheckReqs = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        metricNameSourceAvailabilityCheckRequestTotal,
		Help:        "availability Check requests count partitioned by type (aws/gcp/azure), source status, component and error",
		ConstLabels: prometheus.Labels{"service": "provisioning"},
	},
	[]string{"type", "status", "component", "error"},
)

func IncTotalAvailablilityCheckReqs(provider models.ProviderType, component string, StatusType kafka.StatusType, err error) {
	errString := "false"
	if err != nil {
		errString = "true"
	}
	TotalAvailabilityCheckReqs.WithLabelValues(provider.String(), string(StatusType), component, errString).Inc()
}

func RegisterTotalAvailablilityCheckReqs() {
	prometheus.MustRegister(TotalAvailabilityCheckReqs)
}
