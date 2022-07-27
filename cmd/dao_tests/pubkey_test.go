//go:build integration
// +build integration

package main

import (
	"context"
	"os"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/sqlx"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
)

func Setup(t *testing.T, s string) (dao.PubkeyDao, context.Context, error) {
	err := db.Seed("dao_test")
	if err != nil {
		t.Errorf("Error purging the database: %v", err)
		return nil, nil, err
	}
	ctx := identity.WithTenant(t, context.Background())
	//ctx := context.Background()
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

func TestMain(t *testing.M) {
	exitVal := t.Run()
	os.Exit(exitVal)
}
