package metrics

import "github.com/prometheus/client_golang/prometheus"

func RegisterStatuserMetrics() {
	prometheus.MustRegister(
		TotalSentAvailabilityCheckReqs,
		AvailabilityCheckReqsDuration,
		TotalInvalidAvailabilityCheckReqs,
		RbacAclFetchDuration,
		CacheHits,
	)
}

func RegisterStatsMetrics() {
	prometheus.MustRegister(
		JobQueueSize,
		JobQueueInFlight,
		DbStatsDuration,
		Reservations24hCount,
		Reservations28dCount,
	)
}

func RegisterApiMetrics() {
	prometheus.MustRegister(
		RbacAclFetchDuration,
		CacheHits,
		AvailabilityBatchSendDuration,
	)
}

func RegisterWorkerMetrics() {
	prometheus.MustRegister(
		BackgroundJobDuration,
		ReservationCount,
		RbacAclFetchDuration,
		CacheHits,
	)
}
