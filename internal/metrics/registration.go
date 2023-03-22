package metrics

import "github.com/prometheus/client_golang/prometheus"

func RegisterStatuserMetrics() {
	prometheus.MustRegister(TotalSentAvailabilityCheckReqs, AvailabilityCheckReqsDuration, TotalInvalidAvailabilityCheckReqs)
}

func RegisterApiMetrics() {
	// no metrics
}

func RegisterWorkerMetrics() {
	prometheus.MustRegister(JobQueueSize, JobQueueInFlight, BackgroundJobDuration, ReservationCount)
}
