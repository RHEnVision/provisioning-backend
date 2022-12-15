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
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue/dejq"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	_ "github.com/RHEnVision/provisioning-backend/internal/logging/testing"
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

func initEnvironment() {
	config.Initialize("config/test.env", "../../../config/test.env")
	logging.InitializeStdout()
	ctx := ctxval.WithLogger(context.Background(), &log.Logger)
	dejq.Initialize(ctx, &log.Logger)
	dejq.RegisterJobs(&log.Logger)
	dejq.StartDequeueLoop(ctx, &log.Logger)

	err := db.Initialize(context.Background(), "integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %v", err))
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
	defer dejq.StopDequeueLoop()

	dbDrop()
	dbMigrate()
	exitVal := t.Run()
	dbDrop()
	os.Exit(exitVal)
}
