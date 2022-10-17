package cache

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type noopCache struct{}

func NewNoopCache() *noopCache {
	return &noopCache{}
}

func (_ *noopCache) FindAccountId(_ context.Context, _, _ string) (*models.Account, error) {
	return nil, NotFound
}

func (_ *noopCache) SetAccountId(_ context.Context, _, _ string, _ *models.Account) error {
	// noop
	return nil
}

func (_ *noopCache) FindAppTypeId(_ context.Context) (string, error) {
	return "", NotFound
}

func (_ *noopCache) SetAppTypeId(_ context.Context, _ string) error {
	// noop
	return nil
}
