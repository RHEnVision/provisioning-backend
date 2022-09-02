//go:build !test && !integration
// +build !test,!integration

package config

import (
	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
)

func init() {
	parser.Viper.SetDefault("featureFlags.environment", "development")

	// load configs/local.yaml
	parser.Viper.AddConfigPath("./configs")
	parser.Viper.SetConfigName("local")
	_ = parser.Viper.ReadInConfig()
}
