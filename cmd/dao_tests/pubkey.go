package main

import (
	"context"
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
		ID:        10,
		AccountID: 2,
		Name:      "testkey",
		Body:      "sha-rsa body",
	}}
}

func updatePk() []*models.Pubkey {
	return []*models.Pubkey{{
		ID:        1,
		AccountID: 1,
		Name:      "updated-key",
		Body:      "updated body",
	}}
}

func TearDown(t *testing.T) {
	config.Initialize()

	err := db.Initialize()
	if err != nil {
		t.Errorf("Error initializing database.")
		return
	}

	err = db.Seed("drop_all")
	if err != nil {
		t.Errorf("Error purging the database")
		return
	}

	err = db.Migrate()
	if err != nil {
		t.Errorf("Error running migration")
		return
	}
}

func CreatePubkeyTest(t *testing.T) {
	TearDown(t)
	err := db.Seed("dao_account")
	if err != nil {
		t.Errorf("Error purging the database")
		return
	}

	ctx := context.Background()
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("Create pubkey test had failed.")
		return
	}
	err = pkDao.Create(ctx, createPk()[0])
	if err != nil {
		t.Errorf("Create pubkey test had failed. %s", err)
		return
	}

	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Create pubkey test had failed. %s", err)
		return
	}

	assert.Equal(t, 1, len(pubkeys), "Create pubkey test had failed.")
}

func ListPubkeyTest(t *testing.T) {
	TearDown(t)
	err := db.Seed("dev_small")
	if err != nil {
		t.Errorf("Error purging the database")
		return
	}

	ctx := context.Background()
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("List pubkey test had failed.")
		return
	}
	pubkeys, err := pkDao.List(ctx, 100, 0)
	if err != nil {
		t.Errorf("List pubkey test had failed.")
		return
	}
	assert.Equal(t, 1, len(pubkeys), "List Pubkey error.")
}

func UpdatePubkeyTest(t *testing.T) {
	TearDown(t)
	err := db.Seed("dev_small")
	if err != nil {
		t.Errorf("Error purging the database")
		return
	}

	ctx := context.Background()
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("Update pubkey test had failed.")
		return
	}

	err = pkDao.Update(ctx, updatePk()[0])
	if err != nil {
		t.Errorf("Update pubkey test had failed. %s", err)
		return
	}

	pubkeys, err := pkDao.List(ctx, 10, 0)
	if err != nil {
		t.Errorf("Update pubkey test had failed. %s", err)
		return
	}
	assert.Equal(t, updatePk()[0].Name, pubkeys[0].Name, "Update pubkey test had failed.")
}

func GetPubkeyByIdTest(t *testing.T) {
	TearDown(t)
	err := db.Seed("dev_small")
	if err != nil {
		t.Errorf("Error purging the database")
		return
	}

	ctx := context.Background()
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("Get pubkey test had failed.")
		return
	}

	pubkey, err := pkDao.GetById(ctx, 1)
	if err != nil {
		t.Errorf("Get pubkey test had failed.")
		return
	}
	assert.Equal(t, "lzap-ed25519-2021", pubkey.Name, "Get Pubkey error.")
	assert.Equal(t, "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap", pubkey.Body, "Get Pubkey error.")

}

func DeletePubkeyByIdTest(t *testing.T) {
	TearDown(t)
	err := db.Seed("dev_small")
	if err != nil {
		t.Errorf("Error purging the database")
		return
	}

	ctx := context.Background()
	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		t.Errorf("Delete pubkey test had failed.")
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

func main() {
	testing.Main(
		nil,
		[]testing.InternalTest{
			{"List Pubkeys", ListPubkeyTest},
			{"Create Pubkey", CreatePubkeyTest},
			{"Update Pubkey", UpdatePubkeyTest},
			{"Get Pubkey by ID", GetPubkeyByIdTest},
			{"Delete Pubkey by ID", DeletePubkeyByIdTest},
		},
		nil, nil,
	)
}
