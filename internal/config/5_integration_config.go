//go:build integration && !test
// +build integration,!test

package config

import (
	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
)

func init() {
	parser.Viper.Set("featureFlags.environment", "test")

	// load configs/local.integration.yaml
	parser.Viper.AddConfigPath("./configs")
	// and relative to internal/dao/tests (go test utility changes workdir)
	parser.Viper.AddConfigPath("../../../configs")
	parser.Viper.SetConfigName("local.integration")
	_ = parser.Viper.ReadInConfig()
}
