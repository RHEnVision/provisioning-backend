//go:build integration
// +build integration

// To override application configuration for integration tests, create config/test.env file.
package tests

import (
	"context"
	"os"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"
	_ "github.com/RHEnVision/provisioning-backend/internal/logging/testing"
	"github.com/RHEnVision/provisioning-backend/internal/testing/integration"
)

func TestMain(t *testing.M) {
	ctx := context.Background()
	ctx = integration.InitConfigEnvironment(ctx, "../../../config/test.env")
	config.Application.Cache.Type = "redis"
	integration.InitRedisEnvironment(ctx)
	defer integration.CloseRedisEnvironment(ctx)

	exitVal := t.Run()
	os.Exit(exitVal)
}
