// Copyright Red Hat

package config

import (
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"

	"github.com/spf13/viper"
)

// IdpConfig represents the runtime configuration
type IdpConfig struct {
	Hostname    string
	Auth        bool
	WebPort     int
	MetricsPort int
	Logging     *loggingConfig
	LogLevel    string
	Debug       bool
    Database        *dbConfig
	OpenAPIFilePath string
}

type dbConfig struct {
	User     string
	Password string
	Hostname string
	Port     uint
	Name     string
}

type loggingConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	LogGroup        string
	Region          string
}

var config *IdpConfig

// Init configuration for service
func Init() {
	options := viper.New()
	options.SetDefault("WebPort", 3000)
	options.SetDefault("MetricsPort", 8080)
	options.SetDefault("LogLevel", "INFO")
	options.SetDefault("Auth", true)
	options.SetDefault("Debug", false)
	options.SetDefault("OpenAPIFilePath", "./cmd/spec/openapi.json")
	options.SetDefault("Database", "sqlite")
	options.AutomaticEnv()

	if options.GetBool("Debug") {
		options.Set("LogLevel", "DEBUG")
	}

	kubenv := viper.New()
	kubenv.AutomaticEnv()

	config = &IdpConfig{
		Hostname:        kubenv.GetString("Hostname"),
		Auth:            options.GetBool("Auth"),
		WebPort:         options.GetInt("WebPort"),
		MetricsPort:     options.GetInt("MetricsPort"),
		Debug:           options.GetBool("Debug"),
		LogLevel:        options.GetString("LogLevel"),
		OpenAPIFilePath: options.GetString("OpenAPIFilePath"),
	}

	database := options.GetString("database")

	if database == "pgsql" {
		config.Database = &dbConfig{
			User:     options.GetString("PGSQL_USER"),
			Password: options.GetString("PGSQL_PASSWORD"),
			Hostname: options.GetString("PGSQL_HOSTNAME"),
			Port:     options.GetUint("PGSQL_PORT"),
			Name:     options.GetString("PGSQL_DATABASE"),
		}
	}

	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		config.WebPort = *cfg.PublicPort
		config.MetricsPort = cfg.MetricsPort

		config.Database = &dbConfig{
			User:     cfg.Database.Username,
			Password: cfg.Database.Password,
			Hostname: cfg.Database.Hostname,
			Port:     uint(cfg.Database.Port),
			Name:     cfg.Database.Name,
		}

		config.Logging = &loggingConfig{
			AccessKeyID:     cfg.Logging.Cloudwatch.AccessKeyId,
			SecretAccessKey: cfg.Logging.Cloudwatch.SecretAccessKey,
			LogGroup:        cfg.Logging.Cloudwatch.LogGroup,
			Region:          cfg.Logging.Cloudwatch.Region,
		}
	}
}

// Get returns an initialized IdpConfig
func Get() *IdpConfig {
	return config
}
