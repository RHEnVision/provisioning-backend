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
		viper.Set("restEndpoints.sources.appId", 11)
		viper.Set("featureFlags.environment", "production")
		viper.Set("database.host", cfg.Database.Hostname)
		viper.Set("database.port", cfg.Database.Port)
		viper.Set("database.user", cfg.Database.Username)
		viper.Set("database.password", cfg.Database.Password)
		viper.Set("database.name", cfg.Database.Name)
		viper.Set("prometheus.port", cfg.MetricsPort)
		viper.Set("prometheus.path", cfg.MetricsPath)
	} else {
		viper.SetDefault("featureFlags.environment", "developement")
		viper.SetDefault("restEndpoints.sources.appId", 11)

		viper.AddConfigPath("./configs")
		viper.SetConfigName("defaults")

		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("local")
		err = viper.MergeInConfig()
		if err != nil {
			fmt.Println("Could not read local.yaml", err)
		}
	}
}
