package config

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"regexp"
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
		Region  string
		Key     string
		Secret  string
		Session string
		Group   string
		Stream  string
	}
	AWS struct {
		Region  string
		Key     string
		Secret  string
		Session string
	}
	Prometheus struct {
		Port int
	}
	FeatureFlags struct {
		Environment string
		ExitOnPanic bool
	}
}

var Database = &config.Database
var Prometheus = &config.Prometheus
var Logging = &config.Logging
var Cloudwatch = &config.Cloudwatch
var AWS = &config.AWS
var Features = &config.FeatureFlags

func init() {
	var err error

	viper.SetDefault("logging.level", 1)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.loglevel", 1)
	viper.SetDefault("cloudwatch.enabled", false)
	viper.SetDefault("prometheus.port", 9000)

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		viper.Set("database.host", cfg.Database.Hostname)
		viper.Set("database.port", cfg.Database.Port)
		viper.Set("database.user", cfg.Database.Username)
		viper.Set("database.password", cfg.Database.Password)
		viper.Set("database.name", cfg.Database.Name)
		viper.Set("prometheus.port", cfg.MetricsPort)
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

	if err = Validate(); err != nil {
		panic(err)
	}
}

var validateMissingSecretError = errors.New("config error: Cloudwatch enabled but Region and Key and Secret are not provided")
var validateGroupStreamError = errors.New("config error: Cloudwatch enabled but Group or Stream is blank")
var validateInvalidEnvironmentError = errors.New("config error: Environment must be production or development")

func Validate() error {
	if envMatch, _ := regexp.MatchString(`^(production|development)$`, Features.Environment); !envMatch {
		return fmt.Errorf("%w: %s", validateInvalidEnvironmentError, Features.Environment)
	}

	if Cloudwatch.Enabled {
		if Cloudwatch.Region == "" || Cloudwatch.Key == "" || Cloudwatch.Secret == "" {
			return validateMissingSecretError
		}
		if Cloudwatch.Group == "" || Cloudwatch.Stream == "" {
			return validateGroupStreamError
		}
	}

	return nil
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
	configCopy.Cloudwatch.Key = "*****"
	configCopy.Cloudwatch.Secret = "*****"
	configCopy.Cloudwatch.Session = "*****"
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
