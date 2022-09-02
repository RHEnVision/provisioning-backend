//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func createLukas2021Key() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      "lzap-ed25519-2021",
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap-2021",
	}
}

func createLukas2013Key() *models.Pubkey {
	return &models.Pubkey{
		AccountID: 1,
		Name:      "lzap-rsa-2013",
		Body:      "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dFInru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJDt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDFb2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizAM6pCff3RBslbFxLdOO7cR17 lzap-2013",
	}
}

func setupPubkey(t *testing.T) (dao.PubkeyDao, context.Context) {
	setup()
	ctx := identity.WithTenant(t, context.Background())
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		panic(err)
	}
	return pkDao, ctx
}

func teardownPubkey(_ *testing.T) {
	teardown()
}

func TestCreatePubkey(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	pk := createLukas2013Key()
	err := pkDao.Create(ctx, pk)
	if err != nil {
		t.Errorf("Create pubkey test had failed: %v", err)
		return
	}

	pk2, err := pkDao.GetById(ctx, pk.ID)
	if err != nil {
		t.Errorf("Create pubkey test had failed: %v", err)
		return
	}

	assert.Equal(t, pk.Name, pk2.Name, "Create pubkey test had failed.")
}

func TestCreatePubkeyFingerprintSuccess(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	pk := createLukas2013Key()
	err := pkDao.Create(ctx, pk)
	if err != nil {
		t.Errorf("Create pubkey test had failed: %v", err)
		return
	}

	assert.Equal(t, "SHA256:ENShRe/0uDLSw9c+7tc9PxkD/p4blyB/DTgBSIyTAJY", pk.Fingerprint)
}

func TestCreatePubkeyFingerprintFailure(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	pk := createLukas2021Key()
	err := pkDao.Create(ctx, pk)
	assert.Error(t, err, "key should already exist")
}

func TestListPubkey(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	err := pkDao.Create(ctx, createLukas2013Key())
	pubkeys, err := pkDao.List(ctx, 100, 0)
	if err != nil {
		t.Errorf("List pubkey test had failed: %v", err)
		return
	}
	assert.Equal(t, 2, len(pubkeys), "List Pubkey error.")
}

func TestUpdatePubkey(t *testing.T) {
	updatePk := &models.Pubkey{
		ID:        1,
		AccountID: 1,
		Name:      "avitova-ed25519-2021",
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN avitova",
	}
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	err := pkDao.Create(ctx, createLukas2013Key())
	if err != nil {
		t.Errorf("Create pubkey test failed. %s", err)
		return
	}
	err = pkDao.Update(ctx, updatePk)
	if err != nil {
		t.Errorf("Update pubkey test failed. %s", err)
		return
	}

	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Update pubkey test failed. %s", err)
		return
	}
	assert.Equal(t, updatePk.Name, pubkeys[0].Name, "Update pubkey test had failed.")
}

func TestGetPubkeyById(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	err := pkDao.Create(ctx, createLukas2013Key())
	if err != nil {
		t.Errorf("Delete pubkey test had failed. %s", err)
		return
	}
	pubkey, err := pkDao.GetById(ctx, 1)
	if err != nil {
		t.Errorf("Get pubkey test had failed.")
		return
	}
	assert.Equal(t, "lzap-ed25519-2021", pubkey.Name, "Get Pubkey error: pubkey name does not match.")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap", pubkey.Body, "Get Pubkey error: pubkey body does not match.")

}

func TestDeletePubkeyById(t *testing.T) {
	pkDao, ctx := setupPubkey(t)
	defer teardownPubkey(t)
	err := pkDao.Create(ctx, createLukas2013Key())
	if err != nil {
		t.Errorf("Delete pubkey test had failed. %s", err)
		return
	}
	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Delete pubkey test had failed")
		return
	}
	err = pkDao.Delete(ctx, 1)
	if err != nil {
		t.Errorf("Delete pubkey test had failed")
		return
	}
	pubkeysAfter, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Delete pubkey test had failed")
		return
	}
	assert.Equal(t, len(pubkeys)-1, len(pubkeysAfter), "Delete Pubkey error.")
}
