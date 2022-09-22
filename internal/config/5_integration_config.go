//go:build integration && !test
// +build integration,!test

package config

import (
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
)

func init() {
	// App
	parser.Viper.Set("featureFlags.environment", "test")
	parser.Viper.SetDefault("telemetry.loggerExporter", false)

	// App Cache
	parser.Viper.SetDefault("app.cache.expiration", time.Duration(0))
	parser.Viper.SetDefault("app.cache.cleanupInterval", time.Duration(0))
	parser.Viper.SetDefault("app.cache.typeAppId", false)
	parser.Viper.SetDefault("app.cache.account", false)

	// load configs/local.integration.yaml
	parser.Viper.AddConfigPath("./configs")
	// and relative to internal/dao/tests (go test utility changes workdir)
	parser.Viper.AddConfigPath("../../../configs")
	parser.Viper.SetConfigName("local.integration")
	_ = parser.Viper.ReadInConfig()
}
