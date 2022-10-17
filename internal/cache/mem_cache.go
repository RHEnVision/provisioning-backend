package cache

import (
	"context"
	"sync"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type memCache struct {
	appTypeId      *string
	appTypeIdMutex sync.Mutex
	accountId      *Cache[accountKey, *models.Account]
}

type accountKey struct {
	OrgID         string
	AccountNumber string
}

func NewMemoryCache() *memCache {
	c := memCache{
		accountId: newMemoryCache[accountKey, *models.Account](config.Application.Cache.Memory.CleanupInterval),
	}

	_ = promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "app_cache_account_items",
		Help:        "The total number of cache items for account ID",
		ConstLabels: prometheus.Labels{"service": version.PrometheusLabelName},
	}, func() float64 {
		return float64(c.accountId.Count())
	})

	return &c
}

func (c *memCache) FindAccountId(_ context.Context, OrgID, AccountNumber string) (*models.Account, error) {
	key := accountKey{OrgID: OrgID, AccountNumber: AccountNumber}
	value, ok := c.accountId.Get(key)

	if !ok {
		return nil, NotFound
	}

	return value, nil
}

func (c *memCache) SetAccountId(_ context.Context, OrgID, AccountNumber string, account *models.Account) error {
	key := accountKey{OrgID: OrgID, AccountNumber: AccountNumber}

	c.accountId.Set(key, account, config.Application.Cache.Expiration)
	return nil
}

func (c *memCache) FindAppTypeId(_ context.Context) (string, error) {
	c.appTypeIdMutex.Lock()
	defer c.appTypeIdMutex.Unlock()

	if c.appTypeId == nil {
		return "", NotFound
	}
	return *c.appTypeId, nil
}

func (c *memCache) SetAppTypeId(_ context.Context, value string) error {
	c.appTypeIdMutex.Lock()
	defer c.appTypeIdMutex.Unlock()

	c.appTypeId = &value
	return nil
}
