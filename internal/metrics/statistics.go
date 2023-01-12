package metrics

import (
	"context"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	StatisticsLaunchUsageName = "statistics_launch_usage"
	StatisticsLaunchCountName = "statistics_launch_count"
)

var StatisticsLaunchUsage = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name:        StatisticsLaunchUsageName,
		Help:        "statistics: launch usage",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "api", "statistics": "true"},
	},
	[]string{"type", "provider"},
)

var StatisticsLaunchCount = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:        StatisticsLaunchCountName,
		Help:        "statistics: number of instances launched at once",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName, "component": "api", "statistics": "true"},
		// maximum Value in the wizard is 45
		Buckets: []float64{0.5, 1.5, 2.5, 4.5, 8.5, 16.5, 32.5},
	},
	[]string{"type", "provider"},
)

type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LogRecord struct {
	Name  string   `json:"name"`
	Value []KVPair `json:"pairs"`
}

func statsLogger(ctx context.Context, rec LogRecord) {
	ctxval.Logger(ctx).Info().Interface("statistics", rec).Msgf("Statistics: %+v", rec)
}

func LaunchUsageStats(ctx context.Context, it clients.InstanceTypeName, pt models.ProviderType, count int) {
	StatisticsLaunchUsage.WithLabelValues(it.String(), pt.String()).Inc()
	StatisticsLaunchCount.WithLabelValues(it.String(), pt.String()).Observe(float64(count))
	stats := LogRecord{
		Name: StatisticsLaunchUsageName,
		Value: []KVPair{
			{
				Key:   "type",
				Value: it.String(),
			},
			{
				Key:   "provider",
				Value: pt.String(),
			},
			{
				Key:   "count",
				Value: strconv.Itoa(count),
			},
		},
	}
	statsLogger(ctx, stats)
}

func RegisterStatuserStatistics() {
	// no statistics yet
}

func RegisterApiStatistics() {
	prometheus.MustRegister(StatisticsLaunchUsage, StatisticsLaunchCount)
}
