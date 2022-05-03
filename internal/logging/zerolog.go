package logging

import (
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/clouds/aws"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"os"
	"time"

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

func decorate(l zerolog.Logger) zerolog.Logger {
	return l.With().Timestamp().Str("hostname", hostname).Logger()
}

func InitializeStdout() zerolog.Logger {
	conf := config.GetLoggingConfig()
	zerolog.SetGlobalLevel(zerolog.Level(conf.Level))
	return decorate(log.Output(zerolog.ConsoleWriter{Out: os.Stdout}))
}

func InitializeCloudwatch(logger zerolog.Logger) (zerolog.Logger, func(), error) {
	conf := config.GetLoggingConfig()

	if conf.Cloudwatch {
		log.Debug().Msg("Initializing cloudwatch logger")
		cloudWatchWriter, err := cww.NewWithClient(aws.CWL, 500*time.Millisecond, conf.CWGroup, conf.CWStream)
		if err != nil {
			return logger, nil, fmt.Errorf("cannot initialize cloudwatch: %w", err)
		}

		if conf.Stdout {
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
