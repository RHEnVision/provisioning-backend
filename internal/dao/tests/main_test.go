//go:build integration
// +build integration

// To override application configuration for integration tests, copy local.yaml into this directory.

package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/sqlx"
	"github.com/RHEnVision/provisioning-backend/internal/db"
)

func setup() {
	dbSeed()
}

func teardown() {
	// nothing
}

func initEnvironment() {
	config.Initialize()
	err := db.Initialize("integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database, create configs/local.integration.yaml: %v", err))
	}
}

func dbDrop() {
	err := db.Seed("drop_integration")
	if err != nil {
		panic(err)
	}
}

func dbMigrate() {
	err := db.Migrate("integration")
	if err != nil {
		panic(err)
	}
}

func dbSeed() {
	err := db.Seed("dao_test")
	if err != nil {
		panic(err)
	}
}

func TestMain(t *testing.M) {
	initEnvironment()
	dbDrop()
	dbMigrate()
	exitVal := t.Run()
	dbDrop()
	os.Exit(exitVal)
}
