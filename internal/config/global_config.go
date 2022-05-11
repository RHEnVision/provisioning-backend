package config

import (
	"github.com/rs/zerolog"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type GlobalConfig struct {
	Database struct {
		Host        string        `env:"HOST,required" envDefault:"localhost"`
		Port        uint16        `env:"PORT" envDefault:"5432"`
		Name        string        `env:"NAME"`
		User        string        `env:"USER"`
		Password    string        `env:"PASSWORD,unset"`
		MaxOpenConn int           `env:"MAX_OPEN_CONN" envDefault:"25"`
		MaxIdleConn int           `env:"MAX_IDLE_CONN" envDefault:"20"`
		MaxIdleTime time.Duration `env:"MAX_IDLETIME" envDefault:"5m"`
		MaxLifetime time.Duration `env:"MAX_LIFETIME" envDefault:"10m"`
		LogLevel    int           `env:"LOG_LEVEL" envDefault:"1"`
	} `envPrefix:"DB_"`
	Logging struct {
		Level  int  `env:"LEVEL" envDefault:"1"`
		Stdout bool `env:"STDOUT" envDefault:"false"`
	} `envPrefix:"LOGGING_"`
	Cloudwatch struct {
		Enabled bool   `env:"ENABLED" envDefault:"false"`
		Group   string `env:"GROUP"`
		Stream  string `env:"STREAM"`
	} `envPrefix:"CLOUDWATCH_"`
	AWS struct {
		Region  string `env:"REGION"`
		Key     string `env:"KEY,unset"`
		Secret  string `env:"SECRET,unset"`
		Session string `env:"SESSION,unset"`
	} `envPrefix:"AWS_"`
	FeatureFlags struct {
		Environment string `env:"ENVIRONMENT,required" envDefault:"development"`
	} `envPrefix:"FF_"`
	ExitOnPanic bool `env:"EXIT_ON_PANIC" envDefault:"false"`
}

var config GlobalConfig
var configInitialized bool
var configMutex sync.Mutex

func initializeConfig() {
	// read missing default environmental values from .env.default
	err := godotenv.Load(".env.default")
	if err != nil {
		panic(err)
	}

	// read user-defined overrides from .env
	err = godotenv.Overload()
	if err != nil {
		panic(err)
	}

	// parse variables into the Go type
	err = env.Parse(&config)
	if err != nil {
		panic(err)
	}
	configInitialized = true
}

func GetLoggingConfig() *GlobalConfig {
	configMutex.Lock()
	defer configMutex.Unlock()

	if !configInitialized {
		initializeConfig()
	}
	return &config
}

func IsDevelopment() bool {
	return GetLoggingConfig().FeatureFlags.Environment == "development"
}

func IsProduction() bool {
	return GetLoggingConfig().FeatureFlags.Environment == "production"
}

// DumpConfig writes configuration to a logger. It removes all secrets, however, it is never
// recommended to call this function in production environments.
func DumpConfig(logger zerolog.Logger) {
	if IsProduction() {
		logger.Warn().Msg("dumping configuration in production mode")
	}
	configCopy := config
	configCopy.Database.Password = "*****"
	configCopy.AWS.Key = "*****"
	configCopy.AWS.Secret = "*****"
	configCopy.AWS.Session = "*****"
	logger.Info().Msgf("%+v", configCopy)
}
