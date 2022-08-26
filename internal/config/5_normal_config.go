//go:build !test && !integration
// +build !test,!integration

package config

import (
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func init() {
	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig
		parser.Viper.Set("featureFlags.environment", "production")
		parser.Viper.Set("database.host", cfg.Database.Hostname)
		parser.Viper.Set("database.port", cfg.Database.Port)
		parser.Viper.Set("database.user", cfg.Database.Username)
		parser.Viper.Set("database.password", cfg.Database.Password)
		parser.Viper.Set("database.name", cfg.Database.Name)
		parser.Viper.Set("prometheus.port", cfg.MetricsPort)
		parser.Viper.Set("prometheus.path", cfg.MetricsPath)
		if endpoint, ok := clowder.DependencyEndpoints["sources-api"]["svc"]; ok {
			parser.Viper.Set("restEndpoints.sources.url", fmt.Sprintf("http://%s:%d/api/sources/v3.1", endpoint.Hostname, endpoint.Port))
		}
	} else {
		parser.Viper.SetDefault("featureFlags.environment", "development")

		// load from main subdirectory
		parser.Viper.AddConfigPath("./configs")
		// load from working directory for integration tests
		parser.Viper.AddConfigPath(".")
		parser.Viper.SetConfigName("local")

		err := parser.Viper.ReadInConfig()
		if err != nil {
			fmt.Println("Could not read configs/local.yaml, continuing with defaults")
		}
	}
}
