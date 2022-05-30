package config

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var config struct {
	App struct {
		Name string
		Port int
	}
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
		Path string
	}
	FeatureFlags struct {
		Environment string
		ExitOnPanic bool
	}
}

var Application = &config.App
var Database = &config.Database
var Prometheus = &config.Prometheus
var Logging = &config.Logging
var Cloudwatch = &config.Cloudwatch
var AWS = &config.AWS
var Features = &config.FeatureFlags

func Initialize() {
	err := viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	if err = validate(); err != nil {
		panic(err)
	}
}

var validateMissingSecretError = errors.New("config error: Cloudwatch enabled but Region and Key and Secret are not provided")
var validateGroupStreamError = errors.New("config error: Cloudwatch enabled but Group or Stream is blank")
var validateInvalidEnvironmentError = errors.New("config error: Environment must be production or development")

func validate() error {
	if envMatch, _ := regexp.MatchString(`^(production|development|test)$`, Features.Environment); !envMatch {
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
