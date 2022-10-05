//go:build integration
// +build integration

package tests

import (
	"context"
	"math"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newLukas2021Key() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      factories.GetSequenceName("lzap-ed25519-2021"),
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap-2021",
	}
}

func newLukas2013Key() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      factories.GetSequenceName("lzap-rsa-2013"),
		Body: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dF" +
			"Inru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJ" +
			"Dt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDF" +
			"b2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizA" +
			"M6pCff3RBslbFxLdOO7cR17 lzap-2013",
	}
}

func newLukas2011Key() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      factories.GetSequenceName("lzap-dsa-2011"),
		Body: "ssh-dss AAAAB3NzaC1kc3MAAACBAKqezP3rkK/NcWvMWqoP3qOggGG4QW1vhQJOfyH/l9CbdRxlrcTV9AD5" +
			"BYMcJNn3Ill0iu9d7gSQTZJu2cEWiE8yHJhWOerfPDB4R8BGQlMvbO+8rTplm1Eo3WxtYD0q45Urfh/Ej7HgliTsAYB" +
			"YrQZ0a09auzBlqR3XwH74MlPdAAAAFQClJSTbX6Hp9HzqXyw0P7HeXt0LrwAAAIAogU1yFPDn7xPPUEh16u3ceaZp5w" +
			"H2wDzPEjMHPv+GQd2/yiJB5TX5s9Z5HQax/r3NFhYKNzjyQf1alChS8M0ge9vtPx3oH3Q3NyJGo2wpyYzvDXzP9OHO6" +
			"Vh3PVVOcGL/TlYbFJUeJb8usjtpb4sLmUNuohwifXNAKzkFj/YpswAAAIAf96KvZqMC91JocjY0L09G2kH+v4Ax30VY" +
			"w3iFlA5LgYnbKEBxEvzM+xZ98uRT//Dmn76F5pFIk/QsHpDSHlx5TIuf1pIm6vzuWtUUQUYKTl+ljuft2FY+FfNW4MZ" +
			"ZKx52kr96AOGKKi+U/MAklf+obqf22XFGvNNu4KSjbxqxwg== lzap-2011",
	}
}

func newAnna2021Key() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      factories.GetSequenceName("avitova-nistp256-2021"),
		Body: "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBAaOrIRmMPX84l" +
			"YJ6y3mzH4gBLLCRdeAJX/lsImAn98u3wghha7pD+bp0O9d1iueMVcRpxfnOpxy3hBAoerDjOw= avitova-2021",
	}
}

func setupPubkey(t *testing.T) (dao.PubkeyDao, context.Context) {
	setup()
	ctx := identity.WithTenant(t, context.Background())
	pkDao := dao.GetPubkeyDao(ctx)
	return pkDao, ctx
}

func teardownPubkey(_ *testing.T) {
	teardown()
}

func TestPubkeyCreate(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)

	t.Run("success", func(t *testing.T) {
		pk := newLukas2013Key()
		err := pkDao.Create(ctx, pk)
		require.NoError(t, err)

		pk2, err := pkDao.GetById(ctx, pk.ID)
		require.NoError(t, err)
		assert.Equal(t, pk, pk2)
	})

	t.Run("fingerprint generation", func(t *testing.T) {
		pk := newLukas2011Key()
		err := pkDao.Create(ctx, pk)
		require.NoError(t, err)

		assert.Equal(t, "SHA256:JHqC1mQ9yannTQ32QEi0JfOFbZSPZVq5O68rmpMh7Lo", pk.Fingerprint)
	})

	t.Run("validation error on name", func(t *testing.T) {
		pk := newLukas2021Key()
		pk.Name = ""
		err := pkDao.Create(ctx, pk)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})

	t.Run("validation error on body", func(t *testing.T) {
		pk := newLukas2021Key()
		pk.Body = ""
		err := pkDao.Create(ctx, pk)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})
}

func TestPubkeyList(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)

	t.Run("success", func(t *testing.T) {
		pubkeys, err := pkDao.List(ctx, 1, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pubkeys))
	})

	t.Run("with offset", func(t *testing.T) {
		newKey := newLukas2013Key()
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
	defer teardownPubkey(t)

	t.Run("success", func(t *testing.T) {
		newKey := newLukas2013Key()
		err := pkDao.Create(ctx, newKey)
		defer pkDao.Delete(ctx, newKey.ID)
		require.NoError(t, err)

		updatePk := newAnna2021Key()
		updatePk.ID = newKey.ID
		err = pkDao.Update(ctx, updatePk)
		require.NoError(t, err)

		dbPk, err := pkDao.GetById(ctx, updatePk.ID)
		require.NoError(t, err)
		assert.Equal(t, updatePk, dbPk)
	})

	t.Run("no rows", func(t *testing.T) {
		updatePk := newAnna2021Key()
		updatePk.ID = math.MaxInt64
		err := pkDao.Update(ctx, updatePk)
		require.ErrorIs(t, err, dao.ErrAffectedMismatch)
	})

	t.Run("validation error on name", func(t *testing.T) {
		newKey := newLukas2013Key()
		err := pkDao.Create(ctx, newKey)
		defer pkDao.Delete(ctx, newKey.ID)
		require.NoError(t, err)

		newKey.Name = ""
		err = pkDao.Update(ctx, newKey)
		require.ErrorAs(t, err, &validator.ValidationErrors{})
	})

	t.Run("validation error on body", func(t *testing.T) {
		newKey := newAnna2021Key()
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
	defer teardownPubkey(t)

	t.Run("success", func(t *testing.T) {
		newPk := newLukas2013Key()
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
	defer teardownPubkey(t)

	t.Run("success", func(t *testing.T) {
		pk := newLukas2013Key()
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
