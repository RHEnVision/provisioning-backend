//go:build integration

package code_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/migrations/code"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newContextWithAccount(t *testing.T) context.Context {
	return identity.WithTenant(t, context.Background())
}

func newContext(_ *testing.T) context.Context {
	return context.Background()
}

func TestMigrationUpdateFingerprint(t *testing.T) {
	testCtx := newContext(t)
	ctx := newContextWithAccount(t)

	pkDao := dao.GetPubkeyDao(ctx)
	defer reset()

	t.Run("migrate rsa key", func(t *testing.T) {
		pk := factories.NewPubkeyRSA()
		pk.Type = "test"
		pk.Fingerprint = "12345678901234567890123456789012345678901234"          // 44 chars
		pk.FingerprintLegacy = "12345678901234567890123456789012345678901234567" // 47 chars
		pk.SkipValidation = true
		err := pkDao.Create(ctx, pk)
		require.NotZero(t, pk.ID)
		require.NoError(t, err)

		err = code.UpdateFingerprints(testCtx)
		require.NoError(t, err)

		pk2, err := pkDao.GetById(ctx, pk.ID)
		require.NoError(t, err)
		assert.NotEqual(t, pk, pk2)
		assert.Equal(t, "ENShRe/0uDLSw9c+7tc9PxkD/p4blyB/DTgBSIyTAJY=", pk2.Fingerprint)
		assert.Equal(t, "89:c5:99:b5:33:48:1c:84:be:da:cb:97:45:b0:4a:ee", pk2.FingerprintLegacy)
	})

	t.Run("migrate ed key", func(t *testing.T) {
		pks, err := pkDao.List(ctx, 1, 0) // the key from seed
		require.NoError(t, err)
		pks[0].Type = "test"
		err = pkDao.Update(ctx, pks[0])
		require.NoError(t, err)

		err = code.UpdateFingerprints(testCtx)
		require.NoError(t, err)

		pk2, err := pkDao.GetById(ctx, pks[0].ID)
		require.NoError(t, err)
		assert.Equal(t, pks[0], pk2)
		assert.Equal(t, "gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk=", pk2.Fingerprint)
		assert.Equal(t, "ee:f1:d4:62:99:ab:17:d9:3b:00:66:62:32:b2:55:9e", pk2.FingerprintLegacy)
	})

	t.Run("migrate both rsa and ed keys", func(t *testing.T) {
		pks, err := pkDao.List(ctx, 2, 0)
		require.NoError(t, err)
		for _, pk := range pks {
			pk.Type = "test"
			err = pkDao.Update(ctx, pk)
			require.NoError(t, err)
		}

		err = code.UpdateFingerprints(testCtx)
		require.NoError(t, err)
	})
}
