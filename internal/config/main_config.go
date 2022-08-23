package config

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type proxy struct {
	URL string
}

var config struct {
	App struct {
		Name    string
		Port    int
		Version string
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
		Level    int
		Stdout   bool
		MaxField int
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
	RestEndpoints struct {
		ImageBuilder struct {
			URL      string
			Username string
			Password string
			Proxy    proxy
		}
		Sources struct {
			URL      string
			Username string
			Password string
			Proxy    proxy
		}
		TraceData bool
	}
	Worker struct {
		Queue        string
		Concurrency  int
		HeartbeatSec int
		MaxBeats     int
	}
}

var Application = &config.App
var Database = &config.Database
var Prometheus = &config.Prometheus
var Logging = &config.Logging
var Cloudwatch = &config.Cloudwatch
var AWS = &config.AWS
var Features = &config.FeatureFlags
var RestEndpoints = &config.RestEndpoints
var ImageBuilder = &config.RestEndpoints.ImageBuilder
var Sources = &config.RestEndpoints.Sources
var Worker = &config.Worker

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
	replacement := "****"
	configCopy := config
	configCopy.Database.Password = replacement
	configCopy.Cloudwatch.Key = replacement
	configCopy.Cloudwatch.Secret = replacement
	configCopy.Cloudwatch.Session = replacement
	configCopy.AWS.Key = replacement
	configCopy.AWS.Secret = replacement
	configCopy.AWS.Session = replacement
	configCopy.RestEndpoints.Sources.Password = replacement
	configCopy.RestEndpoints.ImageBuilder.Password = replacement
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
