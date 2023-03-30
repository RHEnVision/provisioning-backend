package metrics

import "github.com/prometheus/client_golang/prometheus"

func RegisterStatuserMetrics() {
	prometheus.MustRegister(TotalSentAvailabilityCheckReqs, AvailabilityCheckReqsDuration, TotalInvalidAvailabilityCheckReqs, CacheHits)
}

func RegisterApiMetrics() {
	prometheus.MustRegister(CacheHits)
}

func RegisterWorkerMetrics() {
	prometheus.MustRegister(JobQueueSize, JobQueueInFlight, BackgroundJobDuration, ReservationCount, CacheHits)
}
