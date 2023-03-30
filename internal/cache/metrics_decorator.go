package cache

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/models"
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
		metrics.IncCacheHit("account", "miss")
		return nil, err
	}

	metrics.IncCacheHit("account", "hit")
	return value, nil
}

// nolint: wrapcheck
func (c *accountMetricsCache) SetAccountId(ctx context.Context, OrgID, AccountNumber string, account *models.Account) error {
	return c.accountId.SetAccountId(ctx, OrgID, AccountNumber, account)
}
