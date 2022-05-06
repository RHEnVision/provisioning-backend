package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"sync"

	"github.com/XMLHexagram/emp"
	"github.com/joho/godotenv"
)

type GlobalConfig struct {
	Logging struct {
		Level  int  `emp:"LEVEL"`
		Stdout bool `emp:"STDOUT"`
	} `emp:"prefix:LOGGING_"`
	Cloudwatch struct {
		Enabled bool   `emp:"ENABLED"`
		Group   string `emp:"GROUP"`
		Stream  string `emp:"STREAM"`
	} `emp:"prefix:CLOUDWATCH_"`
	AWS struct {
		Region  string `emp:"REGION"`
		Key     string `emp:"KEY"`
		Secret  string `emp:"SECRET"`
		Session string `emp:"SESSION"`
	} `emp:"prefix:AWS_"`
	PrintConfig bool `emp:"PRINT_CONFIG"`
	ExitOnPanic bool `emp:"EXIT_ON_PANIC"`
}

var loggingConfig *GlobalConfig
var globalConfigInitialized bool
var globalConfigMutex sync.Mutex

func initializeConfig() {
	// read missing default environmental values from .env.default
	godotenv.Load(".env.default")
	// read user-defined overrides from .env
	godotenv.Overload()
	// parse variables into the Go type
	loggingConfig = new(GlobalConfig)
	parser, err := emp.NewParser(&emp.Config{
		AllowEmpty: true,
	})
	if err != nil {
		panic(err)
	}
	err = parser.Parse(loggingConfig)
	if err != nil {
		panic(err)
	}
	globalConfigInitialized = true
}

func GetLoggingConfig() *GlobalConfig {
	globalConfigMutex.Lock()
	defer globalConfigMutex.Unlock()

	if !globalConfigInitialized {
		initializeConfig()
	}
	return loggingConfig
}

// DumpConfig writes configuration to a logger. Warning: Never dump production configuration into logs!
func DumpConfig(logger zerolog.Logger) {
	msg, err := emp.Marshal(loggingConfig)
	if err != nil {
		log.Error().Msg("Unable to marshall configuration")
	}
	log.Debug().Msg(strings.Join(strings.Split(msg, "\n"), " "))
}
