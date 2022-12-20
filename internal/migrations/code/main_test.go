//go:build integration
// +build integration

package code_test

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
	integration.InitEnvironment(ctx, "../../../config/test.env")
	defer integration.CloseEnvironment(ctx)
	defer integration.DbDrop()

	integration.DbDrop()
	integration.DbMigrate()
	reset()
	exitVal := t.Run()
	os.Exit(exitVal)
}
