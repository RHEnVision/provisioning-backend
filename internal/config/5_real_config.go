//go:build !test
// +build !test

package config

import (
	"fmt"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
)

func init() {
	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig
		viper.Set("featureFlags.environment", "production")
		viper.Set("database.host", cfg.Database.Hostname)
		viper.Set("database.port", cfg.Database.Port)
		viper.Set("database.user", cfg.Database.Username)
		viper.Set("database.password", cfg.Database.Password)
		viper.Set("database.name", cfg.Database.Name)
		viper.Set("prometheus.port", cfg.MetricsPort)
		viper.Set("prometheus.path", cfg.MetricsPath)
		if endpoint, ok := clowder.DependencyEndpoints["sources-api"]["svc"]; ok {
			viper.Set("restEndpoints.sources.url", fmt.Sprintf("http://%s:%d/api/sources/v3.1", endpoint.Hostname, endpoint.Port))
		}
	} else {
		viper.SetDefault("featureFlags.environment", "development")

		// load from main subdirectory
		viper.AddConfigPath("./configs")
		// load from working directory for integration tests
		viper.AddConfigPath(".")
		viper.SetConfigName("local")

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("Could not read configs/local.yaml, continuing with defaults")
		}
	}
}
