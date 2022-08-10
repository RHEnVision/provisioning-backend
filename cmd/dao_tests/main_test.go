//go:build integration
// +build integration

// To override application configuration for integration tests, copy local.yaml into this directory.

package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/zerolog/log"
)

func setup() {
	dbSeed()
}

func teardown() {
	// nothing
}

func initEnvironment() {
	config.Initialize()
	log.Logger = logging.InitializeStdout()

	err := db.Initialize("integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database, create cmd/dao_tests/local.yaml: %v", err))
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
