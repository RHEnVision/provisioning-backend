//go:build integration
// +build integration

// To override application configuration for integration tests, copy local.yaml into this directory.

package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"
	_ "github.com/RHEnVision/provisioning-backend/internal/logging/testing"

	"github.com/RHEnVision/provisioning-backend/internal/config"
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

	err := db.Initialize(context.Background(), "integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database, create configs/local.integration.yaml: %v", err))
	}
}

func dbDrop() {
	err := db.Seed(context.Background(), "drop_integration")
	if err != nil {
		panic(err)
	}
}

func dbMigrate() {
	err := db.Migrate(context.Background(), "integration")
	if err != nil {
		panic(err)
	}
}

func dbSeed() {
	err := db.Seed(context.Background(), "dao_test")
	if err != nil {
		panic(err)
	}
}

func TestMain(t *testing.M) {
	initEnvironment()
	defer db.Close()

	dbDrop()
	dbMigrate()
	exitVal := t.Run()
	dbDrop()
	os.Exit(exitVal)
}
