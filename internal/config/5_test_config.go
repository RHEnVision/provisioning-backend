//go:build test && !integration
// +build test,!integration

package config

import (
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
)

func init() {
	// App
	parser.Viper.Set("featureFlags.environment", "test")

	// App Cache
	parser.Viper.SetDefault("app.cache.expiration", time.Duration(0))
	parser.Viper.SetDefault("app.cache.cleanupInterval", time.Duration(0))
	parser.Viper.SetDefault("app.cache.typeAppId", false)
	parser.Viper.SetDefault("app.cache.account", false)
}
