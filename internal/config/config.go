package config

import (
	"errors"
	"fmt"
	"os"
	"path"
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
		InstancePrefix string `env:"INSTANCE_PREFIX" env-default:"" env-description:"prefix for all VMs names"`
		RbacEnabled    bool   `env:"RBAC_ENABLED" env-default:"false" env-description:"RBAC checking (REST_ENDPOINTS_RBAC_URL must be present)"`
		Notifications  struct {
			Enabled bool `env:"ENABLED" env-default:"false" env-description:"notifications enabled"`
		} `env-prefix:"NOTIFICATIONS_"`
		Cache struct {
			Type       string        `env:"TYPE" env-default:"none" env-description:"application cache (none, redis)"`
			Expiration time.Duration `env:"EXPIRATION" env-default:"10m" env-description:"expiration for Redis application cache (time interval syntax)"`
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
	Stats struct {
		JobQueue             time.Duration `env:"JOBQUEUE_INTERVAL" env-default:"1m" env-description:"how often to pull job queue statistics"`
		ReservationsInterval time.Duration `env:"RESERVATIONS_INTERVAL" env-default:"10m" env-description:"how often to pull reservation statistics"`
	} `env-prefix:"STATS_"`
	Reservation struct {
		CleanupEnabled  bool          `env:"CLEANUP_ENABLED" env-default:"false" env-description:"reservation cleanup enabled"`
		Lifetime        time.Duration `env:"LIFETIME" env-default:"8760h" env-description:"how old reservation should be deleted, default equal to 365 days"`
		CleanupInterval time.Duration `env:"CLEANUP_INTERVAL" env-default:"1h" env-description:"how often to cleanup the reservation"`
	} `env-prefix:"RESERVATION_"`
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
		Stdout   bool   `env:"STDOUT" env-default:"true" env-description:"logger standard output, disabled in clowder by default, stdout is still used if there is no other writer"`
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
		Enabled bool   `env:"ENABLED" env-default:"false" env-description:"cloudwatch logging exporter (enabled in clowder)"`
		Region  string `env:"REGION" env-default:"" env-description:"cloudwatch logging AWS region"`
		Key     string `env:"KEY" env-default:"" env-description:"cloudwatch logging key"`
		Secret  string `env:"SECRET" env-default:"" env-description:"cloudwatch logging secret"`
		Session string `env:"SESSION" env-default:"" env-description:"cloudwatch logging session"`
		Group   string `env:"GROUP" env-default:"" env-description:"cloudwatch logging group"`
		Stream  string `env:"STREAM" env-default:"" env-description:"cloudwatch logging stream"`
	} `env-prefix:"CLOUDWATCH_"`
	AWS struct {
		Key               string        `env:"KEY" env-default:"" env-description:"AWS service account key"`
		Secret            string        `env:"SECRET" env-default:"" env-description:"AWS service account secret"`
		Session           string        `env:"SESSION" env-default:"" env-description:"AWS service account session"`
		DefaultRegion     string        `env:"DEFAULT_REGION" env-default:"us-east-1" env-description:"AWS region when not provided"`
		Logging           bool          `env:"LOGGING" env-default:"false" env-description:"AWS service account logging (verbose)"`
		AvailabilityDelay time.Duration `env:"AVAILABILITY_DELAY" env-default:"1s" env-description:"arbitrary delay between sources availability checks (time interval syntax)"`
		AvailabilityRate  float32       `env:"AVAILABILITY_RATE" env-default:"1.0" env-description:"probability rate for availability checks (0.0 = all skipped, 1.0 = nothing skipped)"`
	} `env-prefix:"AWS_"`
	Azure struct {
		TenantID            string `env:"TENANT_ID" env-default:"" env-description:"Azure service account tenant id"`
		ClientID            string `env:"CLIENT_ID" env-default:"" env-description:"Azure service account client id"`
		ClientSecret        string `env:"CLIENT_SECRET" env-default:"" env-description:"Azure service account client secret"`
		ClientPrincipalID   string `env:"CLIENT_PRINCIPAL_ID" env-default:"" env-description:"Azure Principal ID. It is used for lighthouse delegation. It can be object ID of the service principal or Group object id having the service principal as a member"`
		ClientPrincipalName string `env:"CLIENT_PRINCIPAL_NAME" env-default:"RH HCC" env-description:"Azure display name for the offering principal"`
		DefaultRegion       string `env:"DEFAULT_REGION" env-default:"eastus" env-description:"Azure region when not provided"`
		// SubscriptionID is not used in prod environments - used to fetch instance types
		SubscriptionID    string        `env:"SUBSCRIPTION_ID" env-default:"" env-description:"Azure service account subscription id"`
		AvailabilityDelay time.Duration `env:"AVAILABILITY_DELAY" env-default:"1s" env-description:"arbitrary delay between sources availability checks (time interval syntax)"`
		AvailabilityRate  float32       `env:"AVAILABILITY_RATE" env-default:"1.0" env-description:"probability rate for availability checks (0.0 = all skipped, 1.0 = nothing skipped)"`
	} `env-prefix:"AZURE_"`
	GCP struct {
		ProjectID         string        `env:"PROJECT_ID" env-default:"" env-description:"GCP service account project id"`
		JSON              string        `env:"JSON" env-default:"e30K" env-description:"GCP service account credentials (base64 encoded)"`
		DefaultZone       string        `env:"DEFAULT_ZONE" env-default:"us-east4" env-description:"GCP region when not provided"`
		AvailabilityDelay time.Duration `env:"AVAILABILITY_DELAY" env-default:"1s" env-description:"arbitrary delay between sources availability checks (time interval syntax)"`
		AvailabilityRate  float32       `env:"AVAILABILITY_RATE" env-default:"1.0" env-description:"probability rate for availability checks (0.0 = all skipped, 1.0 = nothing skipped)"`
	} `env-prefix:"GCP_"`
	Prometheus struct {
		Port int    `env:"PORT" env-default:"9000" env-description:"prometheus HTTP port"`
		Path string `env:"PATH" env-default:"/metrics" env-description:"prometheus metrics path"`
	} `env-prefix:"PROMETHEUS_"`
	RestEndpoints struct {
		RBAC struct {
			URL      string `env:"URL" env-default:"" env-description:"RBAC URL"`
			Username string `env:"USERNAME" env-default:"" env-description:"RBAC credentials (dev only)"`
			Password string `env:"PASSWORD" env-default:"" env-description:"RBAC credentials (dev only)"`
			Proxy    proxy  `env-prefix:"PROXY_" env-description:"RBAC HTTP proxy (dev only)"`
		} `env-prefix:"RBAC_"`
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
		Queue        string        `env:"QUEUE" env-default:"memory" env-description:"job worker implementation (memory, redis, sqs, postgres)"`
		PollInterval time.Duration `env:"POLL_INTERVAL" env-default:"5s" env-description:"polling interval (network timeout)"`
		Concurrency  int           `env:"CONCURRENCY" env-default:"33" env-description:"amount of worker polling goroutines (effective concurrency)"`
		Timeout      time.Duration `env:"TIMEOUT" env-default:"30m" env-description:"total timeout for a single job to complete (duration)"`
	} `env-prefix:"WORKER_"`
	Unleash struct {
		Enabled     bool   `env:"ENABLED" env-default:"false" env-description:"unleash service (feature flags)"`
		Environment string `env:"ENVIRONMENT" env-default:"" env-description:"unleash environment"`
		Prefix      string `env:"PREFIX" env-default:"provisioning" env-description:"unleash flag prefix"`
		URL         string `env:"URL" env-default:"http://localhost:4242" env-description:"unleash service URL"`
		Token       string `env:"TOKEN" env-default:"" env-description:"unleash service client access token"`
	} `env-prefix:"UNLEASH_"`
	Sentry struct {
		Dsn string `env:"DSN" env-default:"" env-description:"data source name (empty value disables Sentry)"`
	} `env-prefix:"SENTRY_"`
	Kafka struct {
		Enabled          bool     `env:"ENABLED" env-default:"false" env-description:"kafka service enabled"`
		AuthType         string   `env:"AUTH_TYPE" env-default:"" env-description:"kafka authentication type (MTLS, SASL or empty)"`
		SecurityProtocol string   `env:"PROTOCOL" env-default:"" env-description:"kafka SASL security protocol (PLAINTEXT, SSL, SASL_PLAINTEXT, or SASL_SSL, empty means PLAINTEXT)"`
		TlsSkipVerify    bool     `env:"TLS_SKIP_VERIFY" env-default:"false" env-description:"do not verify TLS server certificate"`
		Brokers          []string `env:"BROKERS" env-default:"localhost:9092" env-description:"kafka hostname:port list of brokers"`
		CACert           string   `env:"CA_CERT" env-default:"" env-description:"kafka TLS CA certificate path (use the OS cert store when blank)"`
		SASL             struct {
			Username      string `env:"USERNAME" env-default:"" env-description:"kafka SASL username"`
			Password      string `env:"PASSWORD" env-default:"" env-description:"kafka SASL password"`
			SaslMechanism string `env:"MECHANISM" env-default:"" env-description:"kafka SASL mechanism (scram-sha-512, scram-sha-256 or plain)"`
		} `env-prefix:"SASL_"`
	} `env-prefix:"KAFKA_"`
}

// Config shortcuts
var (
	Application   = &config.App
	Stats         = &config.Stats
	Reservation   = &config.Reservation
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
	RBAC          = &config.RestEndpoints.RBAC
	Worker        = &config.Worker
	Unleash       = &config.Unleash
	Sentry        = &config.Sentry
	Kafka         = &config.Kafka
)

// Errors
var (
	ErrValidateMissingSecret = errors.New("config error: Cloudwatch enabled but Region or Key or Secret are blank")
	ErrValidateGroupStream   = errors.New("config error: Cloudwatch enabled but Group or Stream is blank")
)

var hostname string

func init() {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown-hostname"
	}
	hostname = h
}

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
		if cfg.Database == nil {
			panic("ERROR: Postgres is required in clowder environment")
		}
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
		if cfg.FeatureFlags == nil {
			panic("ERROR: FeatureFlags is required in clowder environment")
		}
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

				// assumption: SASL credentials are always the same for all nodes in a cluster
				if b.Authtype != nil && *b.Authtype != "" {
					config.Kafka.AuthType = string(*b.Authtype)
				}
				if b.Cacert != nil && *b.Cacert != "" {
					config.Kafka.CACert = *b.Cacert
				}
				if b.Sasl != nil {
					if b.SecurityProtocol != nil && *b.SecurityProtocol != "" {
						config.Kafka.SecurityProtocol = *b.SecurityProtocol
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

		// cloudwatch (is blank in ephemeral)
		cw := cfg.Logging.Cloudwatch
		if present(cw.Region, cw.AccessKeyId, cw.SecretAccessKey, cw.LogGroup) {
			config.Cloudwatch.Enabled = true
			config.Cloudwatch.Key = cw.AccessKeyId
			config.Cloudwatch.Secret = cw.SecretAccessKey
			config.Cloudwatch.Region = cw.Region
			config.Cloudwatch.Group = cw.LogGroup
			config.Cloudwatch.Stream = BinaryName()
		}

		// HTTP proxies are not allowed in clowder environment
		config.RestEndpoints.Sources.Proxy.URL = ""
		config.RestEndpoints.ImageBuilder.Proxy.URL = ""

		// endpoints configuration
		if endpoint, ok := clowder.DependencyEndpoints["rbac"]["service"]; ok {
			config.RestEndpoints.RBAC.URL = fmt.Sprintf("http://%s:%d/api/rbac/v1", endpoint.Hostname, endpoint.Port)
		}
		if endpoint, ok := clowder.DependencyEndpoints["sources-api"]["svc"]; ok {
			config.RestEndpoints.Sources.URL = fmt.Sprintf("http://%s:%d/api/sources/v3.1", endpoint.Hostname, endpoint.Port)
		}
		if endpoint, ok := clowder.DependencyEndpoints["image-builder"]["service"]; ok {
			config.RestEndpoints.ImageBuilder.URL = fmt.Sprintf("http://%s:%d/api/image-builder/v1", endpoint.Hostname, endpoint.Port)
		}
	}

	// validate configuration
	if err := validate(); err != nil {
		panic(err)
	}
}

func BinaryName() string {
	if len(os.Args) < 2 {
		return "unknown"
	}
	return path.Base(os.Args[1])
}

func Hostname() string {
	return hostname
}

func HelpText() (string, error) {
	text, err := cleanenv.GetDescription(&config, ptr.To(""))
	if err != nil {
		return "", fmt.Errorf("cannot generate help text: %w", err)
	}
	return text, nil
}
