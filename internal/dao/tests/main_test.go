//go:build integration
// +build integration

// To override application configuration for integration tests, create config/test.env file.

package tests

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	_ "github.com/RHEnVision/provisioning-backend/internal/logging/testing"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/rs/zerolog/log"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
)

func setup() {
	dbSeed()
}

func teardown() {
	// nothing
}

func initEnvironment() context.Context {
	config.Initialize("config/test.env", "../../../config/test.env")
	logging.InitializeStdout()
	ctx := ctxval.WithLogger(context.Background(), &log.Logger)
	jq.Initialize(ctx, &log.Logger)
	jq.RegisterJobs(&log.Logger)
	jq.StartDequeueLoop(ctx)

	err := db.Initialize(context.Background(), "integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %v", err))
	}
	return ctx
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
	ctx := initEnvironment()
	defer db.Close()
	defer jq.StopDequeueLoop(ctx)

	dbDrop()
	dbMigrate()
	exitVal := t.Run()
	dbDrop()
	os.Exit(exitVal)
}
