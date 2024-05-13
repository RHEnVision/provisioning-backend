//go:build integration

package tests

import (
	"context"
	"database/sql"
	"math"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAccount2() *models.Account {
	return &models.Account{
		OrgID:         "200",
		AccountNumber: sql.NullString{String: "2000", Valid: true},
	}
}

func newAccount3() *models.Account {
	return &models.Account{
		OrgID:         "300",
		AccountNumber: sql.NullString{String: "3000", Valid: true},
	}
}

func newAccountNull() *models.Account {
	return &models.Account{
		OrgID:         "400",
		AccountNumber: sql.NullString{},
	}
}

func setupAccount(t *testing.T) (dao.AccountDao, context.Context) {
	ctx := identity.WithTenant(t, context.Background())
	accDao := dao.GetAccountDao(ctx)
	return accDao, ctx
}

func TestAccountCreate(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		acc := newAccount2()
		err := accDao.Create(ctx, acc)
		require.NoError(t, err)

		account, err := accDao.GetById(ctx, 3)
		require.NoError(t, err)
		assert.Equal(t, acc, account)
	})

	t.Run("with null account number", func(t *testing.T) {
		acc := newAccountNull()
		err := accDao.Create(ctx, acc)
		require.NoError(t, err)

		account, err := accDao.GetByOrgId(ctx, "400")
		require.NoError(t, err)
		assert.Equal(t, acc, account)
	})
}

func TestAccountList(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		acc := newAccount2()
		err := accDao.Create(ctx, acc)
		accounts, err := accDao.List(ctx, 100, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, len(accounts))
		require.Contains(t, accounts, acc)
	})

	t.Run("with offset", func(t *testing.T) {
		a2 := newAccount2()
		_ = accDao.Create(ctx, a2)
		a3 := newAccount3()
		_ = accDao.Create(ctx, a3)
		accounts, err := accDao.List(ctx, 1, 2)
		require.NoError(t, err)
		assert.Equal(t, a2.OrgID, accounts[0].OrgID)
		assert.Equal(t, a2.AccountNumber, accounts[0].AccountNumber)
		accounts, err = accDao.List(ctx, 1, 3)
		require.NoError(t, err)
		assert.Equal(t, a3.OrgID, accounts[0].OrgID)
		assert.Equal(t, a3.AccountNumber, accounts[0].AccountNumber)
	})
}

func TestAccountGetById(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		account, err := accDao.GetById(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, "1", account.OrgID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := accDao.GetById(ctx, math.MaxInt64)
		require.ErrorIs(t, err, dao.ErrNoRows)
	})
}

func TestAccountGetByOrgId(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer reset()

	t.Run("success", func(t *testing.T) {
		account, err := accDao.GetByOrgId(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, int64(1), account.ID)
		assert.Equal(t, "1", account.AccountNumber.String)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := accDao.GetByOrgId(ctx, "0")
		require.ErrorIs(t, err, dao.ErrNoRows)
	})
}

func TestAccountGetOrCreateByIdentity(t *testing.T) {
	accDao, ctx := setupAccount(t)
	defer reset()

	t.Run("new record", func(t *testing.T) {
		account, err := accDao.GetOrCreateByIdentity(ctx, "101", "101")
		require.NoError(t, err)
		account, err = accDao.GetByOrgId(ctx, "101")
		assert.Equal(t, "101", account.OrgID)
	})

	t.Run("already exists by org id", func(t *testing.T) {
		account, err := accDao.GetOrCreateByIdentity(ctx, "1", "0")
		require.NoError(t, err)
		assert.Equal(t, "1", account.OrgID)
		assert.Equal(t, "1", account.AccountNumber.String)
	})

	t.Run("already exists by account number", func(t *testing.T) {
		account, err := accDao.GetOrCreateByIdentity(ctx, "0", "1")
		require.NoError(t, err)
		assert.Equal(t, "1", account.OrgID)
		assert.Equal(t, "1", account.AccountNumber.String)
	})
}
