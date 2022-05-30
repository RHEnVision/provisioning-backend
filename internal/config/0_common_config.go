package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Common initialization must be called first.
func init() {
	viper.SetDefault("app.name", "provisioning")
	viper.SetDefault("app.port", 8000)
	viper.SetDefault("logging.level", 1)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.loglevel", 1)
	viper.SetDefault("cloudwatch.enabled", false)
	viper.SetDefault("prometheus.port", 9000)
	viper.SetDefault("prometheus.path", "/metrics")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
}
