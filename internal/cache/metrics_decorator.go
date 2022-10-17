package cache

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	accountIdHits = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "app_cache_account_hits",
		Help:        "The total number of cache hits for account ID",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
	})
	accountIdMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "app_cache_account_miss",
		Help:        "The total number of cache misses for account ID",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
	})
)

type accountMetricsCache struct {
	accountId AccountIdCache
}

func NewAccountDecorator(accountId AccountIdCache) *accountMetricsCache {
	return &accountMetricsCache{
		accountId: accountId,
	}
}

// nolint: wrapcheck
func (c *accountMetricsCache) FindAccountId(ctx context.Context, OrgID, AccountNumber string) (*models.Account, error) {
	value, err := c.accountId.FindAccountId(ctx, OrgID, AccountNumber)
	if err != nil {
		accountIdMisses.Inc()
		return nil, err
	}

	accountIdHits.Inc()
	return value, nil
}

// nolint: wrapcheck
func (c *accountMetricsCache) SetAccountId(ctx context.Context, OrgID, AccountNumber string, account *models.Account) error {
	return c.accountId.SetAccountId(ctx, OrgID, AccountNumber, account)
}
