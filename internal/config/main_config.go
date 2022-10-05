package config

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/rs/zerolog"
)

type proxy struct {
	URL string
}

// config is the main configuration structure. Values are automatically loaded by
// viper library from local.yaml, however, environmental variables must be defined
// in internal/config/parser/known.go.
var config struct {
	App struct {
		Name           string
		Port           int
		Version        string
		Compression    bool
		InstancePrefix string
		Cache          struct {
			Expiration      time.Duration
			CleanupInterval time.Duration
			AppTypeId       bool
			Account         bool
		}
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
		LogLevel    string
	}
	Logging struct {
		Level    int
		Stdout   bool
		MaxField int
	}
	Telemetry struct {
		Enabled bool
		Jaeger  struct {
			Enabled  bool
			Endpoint string
		}
		Logger struct {
			Enabled bool
		}
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
		Key           string
		Secret        string
		Session       string
		DefaultRegion string
		Logging       bool
	}
	Azure struct {
		TenantID       string
		SubscriptionID string
		ClientID       string
		ClientSecret   string
		DefaultRegion  string
	}
	GCP struct {
		JSON        string
		DefaultZone string
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

// Config shortcuts
var (
	Application   = &config.App
	Database      = &config.Database
	Prometheus    = &config.Prometheus
	Logging       = &config.Logging
	Telemetry     = &config.Telemetry
	Cloudwatch    = &config.Cloudwatch
	AWS           = &config.AWS
	Azure         = &config.Azure
	GCP           = &config.GCP
	Features      = &config.FeatureFlags
	RestEndpoints = &config.RestEndpoints
	ImageBuilder  = &config.RestEndpoints.ImageBuilder
	Sources       = &config.RestEndpoints.Sources
	Worker        = &config.Worker
)

// Errors
var (
	validateMissingSecretError      = errors.New("config error: Cloudwatch enabled but Region and Key and Secret are not provided")
	validateGroupStreamError        = errors.New("config error: Cloudwatch enabled but Group or Stream is blank")
	validateInvalidEnvironmentError = errors.New("config error: Environment must be production or development")
)

func Initialize() {
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
	}

	err := parser.Viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	if err = validate(); err != nil {
		panic(err)
	}
}

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

	slice, err := base64.StdEncoding.DecodeString(config.GCP.JSON)
	config.GCP.JSON = string(slice)
	if err != nil {
		return fmt.Errorf("unable to base64-decode GCP JSON config: %w", err)
	}

	var data json.RawMessage
	err = json.Unmarshal(slice, &data)
	if err != nil {
		return fmt.Errorf("unable to parse GCP JSON config: %w", err)
	}

	return nil
}

func IsDevelopment() bool {
	return Features.Environment == "development"
}

func IsProduction() bool {
	return Features.Environment == "production"
}

func StringToURL(ctx context.Context, urlStr string) *url.URL {
	if urlStr == "" {
		return nil
	}
	urlProxy, err := url.Parse(urlStr)
	if err != nil {
		ctxval.Logger(ctx).Warn().Msgf("Unable to parse proxy URL: '%s', ignoring proxy", urlStr)
		return nil
	}
	return urlProxy
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
	configCopy.Azure.ClientID = replacement
	configCopy.Azure.ClientSecret = replacement
	configCopy.GCP.JSON = replacement
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
