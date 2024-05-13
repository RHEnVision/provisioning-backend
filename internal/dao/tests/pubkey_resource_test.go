//go:build integration

package tests

import (
	"context"
	"math"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPubkeyResourceNoop() *models.PubkeyResource {
	return &models.PubkeyResource{
		Tag:      "tag1",
		PubkeyID: 1,
		Provider: models.ProviderTypeNoop,
		SourceID: "1",
		Handle:   factories.SeqNameWithPrefix("handle"),
		Region:   "us-west-1",
	}
}

func newPubkeyResourceAzure() *models.PubkeyResource {
	return &models.PubkeyResource{
		Tag:      "tag1",
		PubkeyID: 1,
		Provider: models.ProviderTypeAzure,
		Handle:   factories.SeqNameWithPrefix("handle"),
		Region:   "us-east-1",
	}
}

func setupPubkeyResource(t *testing.T) (dao.PubkeyDao, context.Context) {
	ctx := identity.WithTenant(t, context.Background())
	pubkeyDao := dao.GetPubkeyDao(ctx)
	return pubkeyDao, ctx
}

func TestPubkeyResourceCreate(t *testing.T) {
	pubkeyDao, ctx := setupPubkeyResource(t)
	defer reset()

	t.Run("empty", func(t *testing.T) {
		resources, err := pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, 1)
		require.NoError(t, err)
		assert.Empty(t, resources)
	})

	t.Run("success", func(t *testing.T) {
		resource := newPubkeyResourceNoop()
		err := pubkeyDao.UnscopedCreateResource(ctx, resource)
		require.NoError(t, err)

		resources, err := pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, resource.PubkeyID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(resources))
		assert.Equal(t, resource, resources[0])

		// Allows same resource with different region
		resource2 := newPubkeyResourceNoop()
		resource2.Region = "us-east-2"
		err = pubkeyDao.UnscopedCreateResource(ctx, resource2)
		require.NoError(t, err)

		// Allows same resource with different source
		resource3 := newPubkeyResourceNoop()
		resource3.SourceID = "5"
		err = pubkeyDao.UnscopedCreateResource(ctx, resource3)
		require.NoError(t, err)

		// Does not allow identical resource
		resource4 := newPubkeyResourceNoop()
		err = pubkeyDao.UnscopedCreateResource(ctx, resource4)
		require.Error(t, err)

		resources, err = pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, resource.PubkeyID)
		require.NoError(t, err)
		assert.Equal(t, 3, len(resources))
	})
}

func TestPubkeyResourceGetByProviderType(t *testing.T) {
	pubkeyDao, ctx := setupPubkeyResource(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		resource := newPubkeyResourceNoop()
		err := pubkeyDao.UnscopedCreateResource(ctx, resource)
		require.NoError(t, err)

		createdResource, err := pubkeyDao.UnscopedGetResourceBySourceAndRegion(ctx, resource.PubkeyID, resource.SourceID, resource.Region)
		require.NoError(t, err)
		assert.Equal(t, resource, createdResource)
	})

	t.Run("no rows", func(t *testing.T) {
		_, err := pubkeyDao.UnscopedGetResourceBySourceAndRegion(ctx, math.MaxInt64, "1234", "us-east-1")
		require.ErrorIs(t, err, dao.ErrNoRows)
	})
}

func TestPubkeyResourceListByPubkeyId(t *testing.T) {
	pubkeyDao, ctx := setupPubkeyResource(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		resourcesBefore, err := pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, 1)
		require.NoError(t, err)

		noopResource := newPubkeyResourceNoop()
		err = pubkeyDao.UnscopedCreateResource(ctx, noopResource)
		require.NoError(t, err)

		azureResource := newPubkeyResourceAzure()
		err = pubkeyDao.UnscopedCreateResource(ctx, azureResource)
		require.NoError(t, err)

		resourcesAfter, err := pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, 1)
		require.NoError(t, err)

		assert.Equal(t, len(resourcesBefore)+2, len(resourcesAfter))
		require.Contains(t, resourcesAfter, noopResource)
		require.Contains(t, resourcesAfter, azureResource)
	})
}

func TestPubkeyResourceDelete(t *testing.T) {
	pubkeyDao, ctx := setupPubkeyResource(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		resource := newPubkeyResourceNoop()
		err := pubkeyDao.UnscopedCreateResource(ctx, resource)
		require.NoError(t, err)

		resources, err := pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, resource.PubkeyID)
		require.NoError(t, err)
		require.Len(t, resources, 1)

		err = pubkeyDao.UnscopedDeleteResource(ctx, resource.ID)
		require.NoError(t, err)

		resources, err = pubkeyDao.UnscopedListResourcesByPubkeyId(ctx, resource.PubkeyID)
		require.NoError(t, err)
		require.Len(t, resources, 0)
	})

	t.Run("mismatch", func(t *testing.T) {
		err := pubkeyDao.UnscopedDeleteResource(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})
}
