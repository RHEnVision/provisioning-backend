package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/ilyakaznacheev/cleanenv"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

type proxy struct {
	URL string `env:"URL" env-default:"" env-description:"proxy URL (dev only)"`
}

var config struct {
	App struct {
		Port           int    `env:"PORT" env-default:"8000" env-description:"HTTP port of the API service"`
		Compression    bool   `env:"COMPRESSION" env-default:"false" env-description:"HTTP payload compression"`
		InstancePrefix string `env:"INSTANCE_PREFIX" env-default:"" env-description:"prefix for all VMs names"`
		Cache          struct {
			Type       string        `env:"TYPE" env-default:"none" env-description:"application cache (none, memory, redis)"`
			Expiration time.Duration `env:"EXPIRATION" env-default:"1h" env-description:"expiration for both memory and Redis (time interval syntax)"`
			Redis      struct {
				Host     string `env:"HOST" env-default:"localhost" env-description:"redis hostname"`
				Port     int    `env:"PORT" env-default:"6379" env-description:"redis port"`
				User     string `env:"USER" env-default:"" env-description:"redis username"`
				Password string `env:"PASSWORD" env-default:"" env-description:"redis password"`
				DB       int    `env:"DB" env-default:"0" env-description:"redis database number"`
			} `env-prefix:"REDIS_"`
			Memory struct {
				CleanupInterval time.Duration `env:"CLEANUP_INTERVAL" env-default:"5m" env-description:"in-memory expiration interval (time interval syntax)"`
			} `env-prefix:"MEM_"`
		} `env-prefix:"CACHE_"`
	} `env-prefix:"APP_"`
	Database struct {
		Host        string        `env:"HOST" env-default:"localhost" env-description:"main database hostname"`
		Port        uint16        `env:"PORT" env-default:"5432" env-description:"main database port"`
		Name        string        `env:"NAME" env-default:"provisioning" env-description:"main database name"`
		User        string        `env:"USER" env-default:"postgres" env-description:"main database username"`
		Password    string        `env:"PASSWORD" env-default:"" env-description:"main database password"`
		SeedScript  string        `env:"SEED_SCRIPT" env-default:"" env-description:"database seed script (dev only)"`
		MinConn     int32         `env:"MIN_CONN" env-default:"2" env-description:"connection pool minimum size"`
		MaxConn     int32         `env:"MAX_CONN" env-default:"50" env-description:"connection pool maximum size"`
		MaxIdleTime time.Duration `env:"MAX_IDLE_TIME" env-default:"15m" env-description:"connection pool idle time (time interval syntax)"`
		MaxLifetime time.Duration `env:"MAX_LIFETIME" env-default:"2h" env-description:"connection pool total lifetime (time interval syntax)"`
		LogLevel    string        `env:"LOG_LEVEL" env-default:"info" env-description:"logging level of database logs"`
	} `env-prefix:"DATABASE_"`
	Logging struct {
		Level    string `env:"LEVEL" env-default:"info" env-description:"logger level (trace, debug, info, warn, error, fatal, panic)"`
		Stdout   bool   `env:"STDOUT" env-default:"true" env-description:"logger standard output"`
		MaxField int    `env:"MAX_FIELD" env-default:"0" env-description:"logger maximum field length (dev only)"`
	} `env-prefix:"LOGGING_"`
	Telemetry struct {
		Enabled bool `env:"ENABLED" env-default:"false" env-description:"open telemetry collecting"`
		Jaeger  struct {
			Enabled  bool   `env:"ENABLED" env-default:"false" env-description:"open telemetry jaeger exporter"`
			Endpoint string `env:"ENDPOINT" env-default:"http://localhost:14268/api/traces" env-description:"jaeger endpoint"`
		} `env-prefix:"JAEGER_"`
		Logger struct {
			Enabled bool `env:"ENABLED" env-default:"false" env-description:"open telemetry logger output (dev only)"`
		} `env-prefix:"LOGGER_"`
	} `env-prefix:"TELEMETRY_"`
	Cloudwatch struct {
		Enabled bool   `env:"ENABLED" env-default:"false" env-description:"cloudwatch logging exporter"`
		Region  string `env:"REGION" env-default:"" env-description:"cloudwatch logging AWS region"`
		Key     string `env:"KEY" env-default:"" env-description:"cloudwatch logging key"`
		Secret  string `env:"SECRET" env-default:"" env-description:"cloudwatch logging secret"`
		Session string `env:"SESSION" env-default:"" env-description:"cloudwatch logging session"`
		Group   string `env:"GROUP" env-default:"" env-description:"cloudwatch logging group"`
		Stream  string `env:"STREAM" env-default:"" env-description:"cloudwatch logging stream"`
	} `env-prefix:"CLOUDWATCH_"`
	AWS struct {
		Key           string `env:"KEY" env-default:"" env-description:"AWS service account key"`
		Secret        string `env:"SECRET" env-default:"" env-description:"AWS service account secret"`
		Session       string `env:"SESSION" env-default:"" env-description:"AWS service account session"`
		DefaultRegion string `env:"DEFAULT_REGION" env-default:"us-east-1" env-description:"AWS region when not provided"`
		Logging       bool   `env:"LOGGING" env-default:"false" env-description:"AWS service account logging (verbose)"`
	} `env-prefix:"AWS_"`
	Azure struct {
		TenantID       string `env:"TENANT_ID" env-default:"" env-description:"Azure service account tenant id"`
		SubscriptionID string `env:"SUBSCRIPTION_ID" env-default:"" env-description:"Azure service account subscription id"`
		ClientID       string `env:"CLIENT_ID" env-default:"" env-description:"Azure service account client id"`
		ClientSecret   string `env:"CLIENT_SECRET" env-default:"" env-description:"Azure service account client secret"`
		DefaultRegion  string `env:"DEFAULT_REGION" env-default:"eastus" env-description:"Azure region when not provided"`
	} `env-prefix:"AZURE_"`
	GCP struct {
		ProjectID   string `env:"PROJECT_ID" env-default:"" env-description:"GCP service account project id"`
		JSON        string `env:"JSON" env-default:"e30K" env-description:"GCP service account credentials (base64 encoded)"`
		DefaultZone string `env:"DEFAULT_ZONE" env-default:"us-east1" env-description:"GCP region when not provided"`
	} `env-prefix:"GCP_"`
	Prometheus struct {
		Port int    `env:"PORT" env-default:"9000" env-description:"prometheus HTTP port"`
		Path string `env:"PATH" env-default:"/metrics" env-description:"prometheus metrics path"`
	} `env-prefix:"PROMETHEUS_"`
	RestEndpoints struct {
		ImageBuilder struct {
			URL      string `env:"URL" env-default:"" env-description:"image builder URL"`
			Username string `env:"USERNAME" env-default:"" env-description:"image builder credentials (dev only)"`
			Password string `env:"PASSWORD" env-default:"" env-description:"image builder credentials (dev only)"`
			Proxy    proxy  `env-prefix:"PROXY_" env-description:"image builder HTTP proxy (dev only)"`
		} `env-prefix:"IMAGE_BUILDER_"`
		Sources struct {
			URL      string `env:"URL" env-default:"" env-description:"sources URL"`
			Username string `env:"USERNAME" env-default:"" env-description:"sources credentials (dev only)"`
			Password string `env:"PASSWORD" env-default:"" env-description:"sources credentials (dev only)"`
			Proxy    proxy  `env-prefix:"PROXY_" env-description:"sources HTTP proxy (dev only)"`
		} `env-prefix:"SOURCES_"`
		TraceData bool `env:"TRACE_DATA" env-default:"true" env-description:"open telemetry HTTP context pass and trace"`
	} `env-prefix:"REST_ENDPOINTS_"`
	Worker struct {
		Queue       string        `env:"QUEUE" env-default:"memory" env-description:"job worker implementation (memory, redis, sqs, postgres)"`
		Concurrency int           `env:"CONCURRENCY" env-default:"50" env-description:"number of goroutines handling jobs"`
		Heartbeat   time.Duration `env:"HEARTBEAT" env-default:"30s" env-description:"heartbeat interval (time interval syntax)"`
		MaxBeats    int           `env:"MAX_BEATS" env-default:"10" env-description:"maximum amount of heartbeats allowed"`
	} `env-prefix:"WORKER_"`
	Unleash struct {
		Enabled     bool   `env:"ENABLED" env-default:"false" env-description:"unleash service (feature flags)"`
		Environment string `env:"ENVIRONMENT" env-default:"" env-description:"unleash environment"`
		Prefix      string `env:"PREFIX" env-default:"app.provisioning" env-description:"unleash flag prefix"`
		URL         string `env:"URL" env-default:"http://localhost:4242" env-description:"unleash service URL"`
		Token       string `env:"TOKEN" env-default:"" env-description:"unleash service client access token"`
	} `env-prefix:"UNLEASH_"`
	Kafka struct {
		Enabled  bool     `env:"ENABLED" env-default:"false" env-description:"kafka service enabled"`
		Brokers  []string `env:"BROKERS" env-default:"localhost:9092" env-description:"kafka hostname:port list of brokers"`
		AuthType string   `env:"AUTH_TYPE" env-default:"" env-description:"kafka authentication type (mtls, sasl or empty)"`
		CACert   string   `env:"CA_CERT" env-default:"" env-description:"kafka TLS CA certificate path"`
		SASL     struct {
			Username         string `env:"USERNAME" env-default:"" env-description:"kafka SASL username"`
			Password         string `env:"PASSWORD" env-default:"" env-description:"kafka SASL password"`
			SaslMechanism    string `env:"MECHANISM" env-default:"" env-description:"kafka SASL mechanism (scram-sha-512, scram-sha-256 or plain)"`
			SecurityProtocol string `env:"PROTOCOL" env-default:"" env-description:"kafka SASL security protocol"`
		} `env-prefix:"SASL_"`
	} `env-prefix:"KAFKA_"`
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
	RestEndpoints = &config.RestEndpoints
	ImageBuilder  = &config.RestEndpoints.ImageBuilder
	Sources       = &config.RestEndpoints.Sources
	Worker        = &config.Worker
	Unleash       = &config.Unleash
	Kafka         = &config.Kafka
)

// Errors
var (
	validateMissingSecretError = errors.New("config error: Cloudwatch enabled but Region and Key and Secret are not provided")
	validateGroupStreamError   = errors.New("config error: Cloudwatch enabled but Group or Stream is blank")
)

// Initialize loads configuration from provided .env files, the first existing file wins.
func Initialize(configFiles ...string) {
	var loaded bool
	for _, configFile := range configFiles {
		if _, err := os.Stat(configFile); err == nil {
			// if config file exists, load it (also loads environmental variables)
			err := cleanenv.ReadConfig(configFile, &config)
			if err != nil {
				panic(err)
			}
			loaded = true
		}
	}

	if !loaded {
		// otherwise use only environmental variables instead
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			panic(err)
		}
	}

	// override some values when Clowder is present
	if clowder.IsClowderEnabled() {
		cfg := clowder.LoadedConfig

		// database
		config.Database.Host = cfg.Database.Hostname
		config.Database.Port = uint16(cfg.Database.Port)
		config.Database.User = cfg.Database.Username
		config.Database.Password = cfg.Database.Password
		config.Database.Name = cfg.Database.Name

		// prometheus
		config.Prometheus.Port = cfg.MetricsPort
		config.Prometheus.Path = cfg.MetricsPath

		// in-memory cache
		if cfg.InMemoryDb == nil {
			panic("ERROR: Redis is required in clowder environment")
		}
		config.App.Cache.Redis.Host = cfg.InMemoryDb.Hostname
		config.App.Cache.Redis.Port = cfg.InMemoryDb.Port
		if cfg.InMemoryDb.Username != nil {
			config.App.Cache.Redis.User = *cfg.InMemoryDb.Username
		}
		if cfg.InMemoryDb.Password != nil {
			config.App.Cache.Redis.Password = *cfg.InMemoryDb.Password
		}

		// feature flags
		config.Unleash.Enabled = true
		url := fmt.Sprintf("%s://%s:%d/api", cfg.FeatureFlags.Scheme, cfg.FeatureFlags.Hostname, cfg.FeatureFlags.Port)
		config.Unleash.URL = url
		if cfg.FeatureFlags.ClientAccessToken != nil {
			config.Unleash.Token = fmt.Sprintf("Bearer %s", *cfg.FeatureFlags.ClientAccessToken)
		}

		// kafka
		if cfg.Kafka != nil {
			config.Kafka.Enabled = true

			config.Kafka.Brokers = make([]string, len(cfg.Kafka.Brokers))
			for i, b := range cfg.Kafka.Brokers {
				config.Kafka.Brokers[i] = fmt.Sprintf("%s:%d", b.Hostname, *b.Port)

				// assumption: TLS/SASL credentials are always the same for all nodes in a cluster
				if b.Authtype != nil && *b.Authtype != "" {
					config.Kafka.AuthType = string(*b.Authtype)
				}
				if b.Cacert != nil && *b.Cacert != "" {
					config.Kafka.CACert = *b.Cacert
				}
				if b.Sasl != nil {
					if b.Sasl.SecurityProtocol != nil && *b.Sasl.SecurityProtocol != "" {
						config.Kafka.SASL.SecurityProtocol = *b.Sasl.SecurityProtocol
					}
					if b.Sasl.SaslMechanism != nil && *b.Sasl.SaslMechanism != "" {
						config.Kafka.SASL.SaslMechanism = *b.Sasl.SaslMechanism
					}
					if b.Sasl.Username != nil && *b.Sasl.Username != "" {
						config.Kafka.SASL.Username = *b.Sasl.Username
					}
					if b.Sasl.Password != nil && *b.Sasl.Password != "" {
						config.Kafka.SASL.Password = *b.Sasl.Password
					}
				}
			}
		}

		// HTTP proxies are not allowed in clowder environment
		config.RestEndpoints.Sources.Proxy.URL = ""
		config.RestEndpoints.ImageBuilder.Proxy.URL = ""

		// endpoints configuration
		if endpoint, ok := clowder.DependencyEndpoints["sources-api"]["svc"]; ok {
			config.RestEndpoints.Sources.URL = fmt.Sprintf("http://%s:%d/api/sources/v3.1", endpoint.Hostname, endpoint.Port)
		}
		if endpoint, ok := clowder.DependencyEndpoints["image-builder"]["svc"]; ok {
			config.RestEndpoints.Sources.URL = fmt.Sprintf("http://%s:%d/api/image-builder/v1", endpoint.Hostname, endpoint.Port)
		}
	}

	// validate configuration
	if err := validate(); err != nil {
		panic(err)
	}
}

func HelpText() (string, error) {
	text, err := cleanenv.GetDescription(&config, ptr.To(""))
	if err != nil {
		return "", fmt.Errorf("cannot generate help text: %w", err)
	}
	return text, nil
}
