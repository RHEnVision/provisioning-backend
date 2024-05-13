//go:build integration

// To override application configuration for integration tests, create config/test.env file.

package tests

import (
	"context"
	"os"
	"testing"

	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"
	_ "github.com/RHEnVision/provisioning-backend/internal/logging/testing"
	"github.com/RHEnVision/provisioning-backend/internal/testing/integration"
)

// truncate and seed database tables
func reset() {
	integration.DbSeed()
}

func TestMain(t *testing.M) {
	ctx := context.Background()
	ctx = integration.InitConfigEnvironment(ctx, "../../../config/test.env")
	integration.InitDbEnvironment(ctx)
	defer integration.CloseDbEnvironment(ctx)
	defer integration.DbDrop()

	integration.DbDrop()
	integration.DbMigrate()
	reset()
	exitVal := t.Run()
	os.Exit(exitVal)
}
