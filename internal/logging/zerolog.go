package logging

import (
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"github.com/RHEnVision/provisioning-backend/internal/clients/cloudwatchlogs"
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
		return ""
	}

	if utf8.RuneCountInString(str) <= length {
		return str
	}

	return string([]rune(str)[:length]) + "...\""
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
func InitializeStdout() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.Level(config.Logging.Level))
	return decorate(log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Kitchen,
		FormatFieldValue: func(i interface{}) string {
			return truncateText(fmt.Sprintf("%s", i), 40)
		},
	}))
}

func InitializeCloudwatch(logger zerolog.Logger) (zerolog.Logger, func(), error) {
	if config.Cloudwatch.Enabled {
		log.Debug().Msg("Initializing cloudwatch logger")
		cloudWatchWriter, err := cww.NewWithClient(cloudwatchlogs.CWL, 500*time.Millisecond, config.Cloudwatch.Group, config.Cloudwatch.Stream)
		if err != nil {
			return logger, nil, fmt.Errorf("cannot initialize cloudwatch: %w", err)
		}

		if config.Logging.Stdout {
			// stdout and cloudwatch
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			newLogger := decorate(zerolog.New(zerolog.MultiLevelWriter(consoleWriter, cloudWatchWriter)))
			return newLogger, cloudWatchWriter.Close, nil
		} else {
			// only cloudwatch (production mode)
			newLogger := decorate(zerolog.New(cloudWatchWriter))
			return newLogger, cloudWatchWriter.Close, nil
		}
	}

	return logger, func() {}, nil
}

func DumpConfigForDevelopment() {
	if config.IsDevelopment() {
		config.DumpConfig(log.Logger)
	}
}
