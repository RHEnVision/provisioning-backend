package logging

import (
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	cww "github.com/lzap/cloudwatchwriter2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var hostname string

func init() {
	h, err := os.Hostname()
	if err != nil {
		h = "unknown-hostname"
	}
	hostname = h
}

func truncateText(str string, length int) string {
	if length <= 0 {
		return str
	}

	if utf8.RuneCountInString(str) <= length {
		return str
	}

	trimmed := []rune(str)[:length]

	if trimmed[0] == '"' {
		return string(trimmed) + "...\""
	} else {
		return string(trimmed) + "..."
	}
}

func decorate(l zerolog.Logger) zerolog.Logger {
	logger := l.With().Timestamp().
		Str("hostname", hostname)

	if version.BuildCommit != "" {
		logger = logger.Str("version", version.BuildCommit)
	}

	return logger.Logger()
}

// InitializeStdout initializes logging to standard output with human-friendly output.
// It is used before CloudWatch logging output is initialized, or in unit and integration tests.
func InitializeStdout() {
	level, err := zerolog.ParseLevel(config.Logging.Level)
	if err != nil {
		panic(fmt.Errorf("cannot parse log level '%s': %w", config.Logging.Level, err))
	}
	zerolog.SetGlobalLevel(level)

	log.Logger = decorate(log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    config.InEphemeralClowder(),
		TimeFormat: time.Kitchen,
		FormatFieldValue: func(i interface{}) string {
			return truncateText(fmt.Sprintf("%s", i), config.Logging.MaxField)
		},
	}))
}

func InitializeCloudwatch(logger zerolog.Logger) (zerolog.Logger, func(), error) {
	if config.Cloudwatch.Enabled {
		log.Debug().Msgf("Initializing cloudwatch logger key %s group %s stream %s region %s",
			config.Cloudwatch.Key, config.Cloudwatch.Group, config.Cloudwatch.Stream, config.Cloudwatch.Region)

		cwClient := newCloudwatchClient()
		cloudWatchWriter, err := cww.NewWithClient(cwClient, 500*time.Millisecond, config.Cloudwatch.Group, config.Cloudwatch.Stream)
		if err != nil {
			return logger, nil, fmt.Errorf("cannot initialize cloudwatch: %w", err)
		}

		if !config.InClowder() && config.Logging.Stdout {
			// stdout and cloudwatch during local development
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			newLogger := decorate(zerolog.New(zerolog.MultiLevelWriter(consoleWriter, cloudWatchWriter)))
			return newLogger, cloudWatchWriter.Close, nil
		} else {
			// only cloudwatch (production mode)
			newLogger := decorate(zerolog.New(cloudWatchWriter))
			return newLogger, cloudWatchWriter.Close, nil
		}
	}

	log.Trace().Msg("Cloudwatch not enabled")
	return logger, func() {}, nil
}

func DumpConfigForDevelopment() {
	// safe dump of app config
	if !config.InProdClowder() {
		config.DumpClowder(log.Logger)
		config.DumpConfig(log.Logger)
	}
}
