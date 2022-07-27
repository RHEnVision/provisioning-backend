//go:build integration
// +build integration

package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/sqlx"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/stretchr/testify/assert"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
)

func createPk() []*models.Pubkey {
	return []*models.Pubkey{{
		ID:        1,
		AccountID: 1,
		Name:      "lzap-ed25519-2021",
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
	}}
}

func Setup(t *testing.T, s string) (dao.PubkeyDao, context.Context, error) {
	err := db.Seed("dao_test")
	if err != nil {
		t.Errorf("Error purging the database: %v", err)
		return nil, nil, err
	}
	ctx := context.Background()
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("%s test had failed: %v", s, err)
		return nil, nil, err
	}
	return pkDao, ctx, nil
}

func CleanUpDatabase(t *testing.T) {
	config.Initialize()

	err := db.Initialize()
	if err != nil {
		t.Errorf("Error initializing database: %v", err)
		return
	}

	err = db.Seed("drop_all")
	if err != nil {
		t.Errorf("Error purging the database: %v", err)
		return
	}

	err = db.Migrate()
	if err != nil {
		t.Errorf("Error running migration: %v", err)
		return
	}
}

func TestSimple(t *testing.T) {
	CleanUpDatabase(t)
	_, _, err := Setup(t, "Simple test")
	if err != nil {
		t.Errorf("Database setup had failed: %v", err)
		return
	}
	assert.Nil(t, nil)
}

func TestCreatePubkey(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "Create pubkey")
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Database setup had failed: %v", err))
		return
	}
	err = pkDao.Create(ctx, createPk()[0])
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Create pubkey test had failed: %v", err))
		return
	}

	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Create pubkey test had failed: %v", err))
		return
	}

	assert.Equal(t, 1, len(pubkeys), "Create pubkey test had failed.")
}

func TestListPubkey(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "List pubkey")
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Database setup had failed: %v", err))
		return
	}
	err = pkDao.Create(ctx, createPk()[0])
	pubkeys, err := pkDao.List(ctx, 100, 0)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("List pubkey test had failed: %v", err))
		return
	}
	assert.Equal(t, 1, len(pubkeys), "List Pubkey error.")
}

func TestUpdatePubkey(t *testing.T) {
	CleanUpDatabase(t)
	updatePk := []*models.Pubkey{{
		ID:        1,
		AccountID: 1,
		Name:      "avitova-ed25519-2021",
		Body:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN avitova",
	}}
	pkDao, ctx, err := Setup(t, "Update pubkey")
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Database setup had failed. %s", err))
		return
	}
	err = pkDao.Create(ctx, createPk()[0])
	err = pkDao.Update(ctx, updatePk[0])
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Update pubkey test had failed. %s", err))
		return
	}

	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Update pubkey test had failed. %s", err))
		return
	}
	assert.Equal(t, updatePk[0].Name, pubkeys[0].Name, "Update pubkey test had failed.")
}

func TestGetPubkeyById(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "Get pubkey")
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Database setup had failed. %s", err))
		return
	}
	err = pkDao.Create(ctx, createPk()[0])
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Delete pubkey test had failed. %s", err))
		return
	}
	pubkey, err := pkDao.GetById(ctx, 1)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Get pubkey test had failed."))
		return
	}
	assert.Equal(t, "lzap-ed25519-2021", pubkey.Name, "Get Pubkey error: pubkey name does not match.")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap", pubkey.Body, "Get Pubkey error: pubkey body does not match.")

}

func TestDeletePubkeyById(t *testing.T) {
	CleanUpDatabase(t)
	pkDao, ctx, err := Setup(t, "Delete pubkey")
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Database setup had failed"))
		return
	}
	err = pkDao.Create(ctx, createPk()[0])
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Delete pubkey test had failed. %s", err))
		return
	}
	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Delete pubkey test had failed"))
		return
	}
	err = pkDao.Delete(ctx, 1)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Delete pubkey test had failed"))
		return
	}
	pubkeysAfter, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		assert.Nil(t, err, fmt.Sprintf("Delete pubkey test had failed"))
		return
	}
	assert.Equal(t, len(pubkeys)-1, len(pubkeysAfter), "Delete Pubkey error.")
}
