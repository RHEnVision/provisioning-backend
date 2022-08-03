package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients/cloudwatchlogs"
	"github.com/RHEnVision/provisioning-backend/internal/config"

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

func truncateText(s string, max int) string {
	if max > len(s) {
		return s
	}
	return s[:strings.LastIndex(s[:max], " ")] + "...\""
}

func decorate(l zerolog.Logger) zerolog.Logger {
	return l.With().Timestamp().Str("hostname", hostname).Logger()
}

func InitializeStdout() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.Level(config.Logging.Level))
	return decorate(log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.Kitchen,
		NoColor:    true,
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
