//go:build integration
// +build integration

package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	value := models.Account{
		ID:            42,
		OrgID:         "442",
		AccountNumber: sql.NullString{},
	}
	err := cache.Set(context.Background(), "42", &value)
	require.NoError(t, err)

	result := models.Account{}
	err = cache.Find(context.Background(), "42", &result)
	require.NoError(t, err)

	require.Equal(t, value, result)
}

func TestAzureTenantId(t *testing.T) {
	var value clients.AzureTenantId = "12345"
	err := cache.Set(context.Background(), "42", &value)
	require.NoError(t, err)

	var result clients.AzureTenantId
	err = cache.Find(context.Background(), "42", &result)
	require.NoError(t, err)

	require.Equal(t, value, result)
}

func TestPrefix(t *testing.T) {
	value1 := models.Account{ID: 1}
	err := cache.Set(context.Background(), "1", &value1)
	require.NoError(t, err)

	// must not overwrite value1
	var value2 clients.AzureTenantId = "2"
	err = cache.Set(context.Background(), "1", &value2)
	require.NoError(t, err)

	result := models.Account{}
	err = cache.Find(context.Background(), "1", &result)
	require.NoError(t, err)

	require.Equal(t, value1, result)
}
