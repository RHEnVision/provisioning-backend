package metrics

import "github.com/prometheus/client_golang/prometheus"

func RegisterStatuserMetrics() {
	prometheus.MustRegister(
		TotalSentAvailabilityCheckReqs,
		AvailabilityCheckReqsDuration,
		TotalInvalidAvailabilityCheckReqs,
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
		CacheHits,
	)
}

func RegisterWorkerMetrics() {
	prometheus.MustRegister(
		BackgroundJobDuration,
		ReservationCount,
		CacheHits,
	)
}
