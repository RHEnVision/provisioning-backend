// Package cache provides application and HTTP response cache.
package cache

import (
	"context"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/rs/zerolog/log"
)

var NotFound = errors.New("not found in cache")

type AccountIdCache interface {
	FindAccountId(ctx context.Context, OrgID, AccountNumber string) (*models.Account, error)
	SetAccountId(ctx context.Context, OrgID, AccountNumber string, account *models.Account) error
}

type AppTypeIdCache interface {
	FindAppTypeId(ctx context.Context) (string, error)
	SetAppTypeId(ctx context.Context, newAppTypeId string) error
}

var (
	accountIdCache AccountIdCache = NewNoopCache()
	appTypeIdCache AppTypeIdCache = NewNoopCache()
)

func Initialize() {
	if config.Application.Cache.Type == "memory" {
		log.Logger.Info().Msg("Initializing memory application cache")
		appTypeIdCache = NewMemoryCache()
		accountIdCache = NewAccountDecorator(NewMemoryCache())
	} else if config.Application.Cache.Type == "redis" {
		log.Logger.Info().Msg("Initializing redis application cache")
		appTypeIdCache = NewMemoryCache()
		accountIdCache = NewAccountDecorator(NewRedisCache())
	} else {
		log.Logger.Info().Msg("No application cache in use")
	}
}

// nolint: wrapcheck
func FindAccountId(ctx context.Context, OrgID, AccountNumber string) (*models.Account, error) {
	return accountIdCache.FindAccountId(ctx, OrgID, AccountNumber)
}

// nolint: wrapcheck
func SetAccountId(ctx context.Context, OrgID, AccountNumber string, account *models.Account) error {
	return accountIdCache.SetAccountId(ctx, OrgID, AccountNumber, account)
}

// nolint: wrapcheck
func FindAppTypeId(ctx context.Context) (string, error) {
	return appTypeIdCache.FindAppTypeId(ctx)
}

// nolint: wrapcheck
func SetAppTypeId(ctx context.Context, value string) error {
	return appTypeIdCache.SetAppTypeId(ctx, value)
}
