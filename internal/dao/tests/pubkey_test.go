//go:build integration

package tests

import (
	"context"
	"math"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPubkey(t *testing.T) (dao.PubkeyDao, context.Context) {
	ctx := identity.WithTenant(t, context.Background())
	pkDao := dao.GetPubkeyDao(ctx)
	return pkDao, ctx
}

func TestPubkeyCreate(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		pk := factories.NewPubkeyRSA()
		err := pkDao.Create(ctx, pk)
		require.NoError(t, err)

		pk2, err := pkDao.GetById(ctx, pk.ID)
		require.NoError(t, err)
		assert.Equal(t, pk, pk2)
	})

	t.Run("fingerprint generation of unsupported key", func(t *testing.T) {
		pk := factories.NewPubkeyDSS()
		err := pkDao.Create(ctx, pk)
		require.ErrorContains(t, err, "x509: unsupported public key type")
	})

	t.Run("validation error on name", func(t *testing.T) {
		pk := factories.NewPubkeyED25519()
		pk.Name = ""
		err := pkDao.Create(ctx, pk)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})

	t.Run("validation error on body", func(t *testing.T) {
		pk := factories.NewPubkeyED25519()
		pk.Body = ""
		err := pkDao.Create(ctx, pk)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})
}

func TestPubkeyList(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		pubkeys, err := pkDao.List(ctx, 1, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pubkeys))
	})

	t.Run("with offset", func(t *testing.T) {
		newKey := factories.NewPubkeyRSA()
		err := pkDao.Create(ctx, newKey)
		require.NoError(t, err)

		pubkeys, err := pkDao.List(ctx, 1, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pubkeys))

		pubkeys, err = pkDao.List(ctx, 1, 1)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pubkeys))
		require.Contains(t, pubkeys, newKey)
	})
}

func TestPubkeyUpdate(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		newKey := factories.NewPubkeyRSA()
		err := pkDao.Create(ctx, newKey)
		defer pkDao.Delete(ctx, newKey.ID)
		require.NoError(t, err)

		updatePk := factories.NewPubkeyECDSA()
		updatePk.ID = newKey.ID
		err = pkDao.Update(ctx, updatePk)
		require.NoError(t, err)

		dbPk, err := pkDao.GetById(ctx, updatePk.ID)
		require.NoError(t, err)
		assert.Equal(t, updatePk, dbPk)
	})

	t.Run("no rows", func(t *testing.T) {
		updatePk := factories.NewPubkeyECDSA()
		updatePk.ID = math.MaxInt64
		err := pkDao.Update(ctx, updatePk)
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})

	t.Run("validation error on name", func(t *testing.T) {
		newKey := factories.NewPubkeyRSA()
		err := pkDao.Create(ctx, newKey)
		defer pkDao.Delete(ctx, newKey.ID)
		require.NoError(t, err)

		newKey.Name = ""
		err = pkDao.Update(ctx, newKey)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})

	t.Run("validation error on body", func(t *testing.T) {
		newKey := factories.NewPubkeyECDSA()
		defer pkDao.Delete(ctx, newKey.ID)
		err := pkDao.Create(ctx, newKey)
		require.NoError(t, err)

		newKey.Body = ""
		err = pkDao.Update(ctx, newKey)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})
}

func TestPubkeyGetById(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		newPk := factories.NewPubkeyRSA()
		err := pkDao.Create(ctx, newPk)
		require.NoError(t, err)

		dbPk, err := pkDao.GetById(ctx, newPk.ID)
		require.NoError(t, err)
		assert.Equal(t, newPk, dbPk)
	})

	t.Run("no rows", func(t *testing.T) {
		_, err := pkDao.GetById(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrNoRows)
	})
}

func TestPubkeyDeleteById(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		pk := factories.NewPubkeyRSA()
		err := pkDao.Create(ctx, pk)
		require.NoError(t, err)

		err = pkDao.Delete(ctx, pk.ID)
		require.NoError(t, err)

		_, err = pkDao.GetById(ctx, pk.ID)
		require.ErrorIs(t, err, dao.ErrNoRows)
	})

	t.Run("mismatch", func(t *testing.T) {
		err := pkDao.Delete(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})
}
