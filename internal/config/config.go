package config

import (
	"fmt"
	"github.com/rs/zerolog"
	"strings"
	"time"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
)

var config struct {
	Database struct {
		Host        string
		Port        uint16
		Name        string
		User        string
		Password    string
		SeedScript  string
		MaxIdleTime time.Duration
		MaxLifetime time.Duration
		MaxOpenConn int
		MaxIdleConn int
		LogLevel    int
	}
	Logging struct {
		Level  int
		Stdout bool
	}
	Cloudwatch struct {
		Enabled bool
		Group   string
		Stream  string
	}
	FeatureFlags struct {
		Environment string
		ExitOnPanic bool
	}
}

var Database = &config.Database
var Logging = &config.Logging
var Cloudwatch = &config.Cloudwatch
var Features = &config.FeatureFlags

func init() {
	var err error

	viper.SetDefault("logging.level", 1)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.loglevel", 1)
	viper.SetDefault("cloudwatch.enabled", false)

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		viper.Set("database.host", cfg.Database.Hostname)
		viper.Set("database.port", cfg.Database.Port)
		viper.Set("database.user", cfg.Database.Username)
		viper.Set("database.password", cfg.Database.Password)
		viper.Set("database.name", cfg.Database.Name)
	} else {
		viper.AddConfigPath("./configs")
		viper.SetConfigName("defaults")
		err = viper.ReadInConfig()
		if err != nil {
			panic(err)
		}
		viper.SetConfigName("local")
		err = viper.MergeInConfig()
		if err != nil {
			fmt.Println("Could not read local.yaml", err)
		}
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func IsDevelopment() bool {
	return Features.Environment == "development"
}

func IsProduction() bool {
	return Features.Environment == "production"
}

// DumpConfig writes configuration to a logger. It removes all secrets, however, it is never
// recommended to call this function in production environments.
func DumpConfig(logger zerolog.Logger) {
	if IsProduction() {
		logger.Warn().Msg("Dumping configuration in production mode!")
	}
	configCopy := config
	configCopy.Database.Password = "*****"
	// configCopy.AWS.Key = "*****"
	// configCopy.AWS.Secret = "*****"
	// configCopy.AWS.Session = "*****"
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
