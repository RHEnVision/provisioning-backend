package cache

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/cache/memcache"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	accountCache = memcache.New[AccountKey, *models.Account](config.Application.Cache.CleanupInterval)
	metricHit    = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "app_cache_account_hits",
		Help:        "The total number of cache hits for account ID",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
	})
	metricMiss = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "app_cache_account_miss",
		Help:        "The total number of cache misses for account ID",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
	})
	_ = promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "app_cache_account_items",
		Help:        "The total number of cache items for account ID",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
	}, func() float64 {
		return float64(accountCache.Count())
	})
)

type AccountKey struct {
	OrgID         string
	AccountNumber string
}

func FindAccountId(ctx context.Context, key AccountKey) (*models.Account, bool) {
	if !config.Application.Cache.Account {
		return nil, false
	}

	if value, ok := accountCache.Get(key); ok {
		metricHit.Inc()
		return value, ok
	}

	metricMiss.Inc()
	return nil, false
}

func SetAccountId(_ context.Context, key AccountKey, account *models.Account) {
	if !config.Application.Cache.Account {
		return
	}

	accountCache.Set(key, account, config.Application.Cache.Expiration)
}
