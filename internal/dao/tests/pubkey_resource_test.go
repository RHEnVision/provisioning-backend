//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPubkeyResourceNoop() *models.PubkeyResource {
	return &models.PubkeyResource{
		Tag:      "tag1",
		PubkeyID: 1,
		Provider: models.ProviderTypeNoop,
		Handle:   factories.GetSequenceName("handle"),
	}
}

func createPubkeyResourceAzure() *models.PubkeyResource {
	return &models.PubkeyResource{
		Tag:      "tag1",
		PubkeyID: 1,
		Provider: models.ProviderTypeAzure,
		Handle:   factories.GetSequenceName("handle"),
	}
}

func setupPubkeyResource(t *testing.T) (dao.PubkeyResourceDao, context.Context) {
	setup()
	ctx := identity.WithTenant(t, context.Background())
	resourceDao, err := dao.GetPubkeyResourceDao(ctx)
	if err != nil {
		panic(err)
	}
	return resourceDao, ctx
}

func teardownPubkeyResource(_ *testing.T) {
	teardown()
}

func TestCreateResource(t *testing.T) {
	resourceDao, ctx := setupPubkeyResource(t)
	defer teardownPubkeyResource(t)
	resource := createPubkeyResourceNoop()
	err := resourceDao.UnscopedCreate(ctx, resource)
	require.NoError(t, err)
	resources, err := resourceDao.UnscopedListByPubkeyId(ctx, resource.PubkeyID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(resources))
	assert.Equal(t, resource, resources[0])
}

func TestGetResourceByProviderType(t *testing.T) {
	resourceDao, ctx := setupPubkeyResource(t)
	defer teardownPubkeyResource(t)
	resource := createPubkeyResourceNoop()
	err := resourceDao.UnscopedCreate(ctx, resource)
	require.NoError(t, err)
	createdResource, err := resourceDao.UnscopedGetResourceByProviderType(ctx, resource.PubkeyID, resource.Provider)
	require.NoError(t, err)

	assert.Equal(t, resource, createdResource)
}

func TestListByPubkeyIdResource(t *testing.T) {
	resourceDao, ctx := setupPubkeyResource(t)
	defer teardownPubkeyResource(t)
	pkId := int64(1)
	resourcesBefore, err := resourceDao.UnscopedListByPubkeyId(ctx, pkId)
	require.NoError(t, err)
	noopResource := createPubkeyResourceNoop()
	err = resourceDao.UnscopedCreate(ctx, noopResource)
	require.NoError(t, err)
	azureResource := createPubkeyResourceAzure()
	err = resourceDao.UnscopedCreate(ctx, azureResource)
	require.NoError(t, err)
	resourcesAfter, err := resourceDao.UnscopedListByPubkeyId(ctx, pkId)
	require.NoError(t, err)

	assert.Equal(t, len(resourcesBefore)+2, len(resourcesAfter))
	require.Contains(t, resourcesAfter, noopResource)
	require.Contains(t, resourcesAfter, azureResource)
}

func TestDeleteResource(t *testing.T) {
	resourceDao, ctx := setupPubkeyResource(t)
	defer teardownPubkeyResource(t)
	resource := createPubkeyResourceNoop()
	err := resourceDao.UnscopedCreate(ctx, resource)
	require.NoError(t, err)
	resources, err := resourceDao.UnscopedListByPubkeyId(ctx, resource.PubkeyID)
	require.NoError(t, err)

	assert.Equal(t, 1, len(resources))

	err = resourceDao.UnscopedDelete(ctx, resource.ID)
	require.NoError(t, err)
	resources, err = resourceDao.UnscopedListByPubkeyId(ctx, resource.PubkeyID)
	require.NoError(t, err)

	assert.Equal(t, 0, len(resources))
}
