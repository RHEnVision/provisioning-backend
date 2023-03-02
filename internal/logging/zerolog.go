package logging

import (
	"fmt"
	"io"
	"os"
	"time"
	"unicode/utf8"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	cww "github.com/lzap/cloudwatchwriter2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

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
		Str("hostname", config.Hostname())

	if version.BuildCommit != "" {
		logger = logger.Str("version", version.BuildCommit)
	}

	return logger.Logger()
}

func configureZerolog() {
	level, err := zerolog.ParseLevel(config.Logging.Level)
	if err != nil {
		panic(fmt.Errorf("cannot parse log level '%s': %w", config.Logging.Level, err))
	}
	zerolog.SetGlobalLevel(level)
	//nolint:reassign
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func stdoutWriter(truncate bool) io.Writer {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    config.InEphemeralClowder(),
		TimeFormat: time.Kitchen,
	}
	if truncate {
		writer.FormatFieldValue = func(i interface{}) string {
			return truncateText(fmt.Sprintf("%s", i), config.Logging.MaxField)
		}
	}
	return writer
}

func cloudwatchWriter() (io.Writer, func(), error) {
	log.Debug().Msgf("Initializing cloudwatch logger key %s group %s stream %s region %s",
		config.Cloudwatch.Key, config.Cloudwatch.Group, config.Cloudwatch.Stream, config.Cloudwatch.Region)

	cwClient := newCloudwatchClient()
	cloudWatchWriter, err := cww.NewWithClient(cwClient, 500*time.Millisecond, config.Cloudwatch.Group, config.Cloudwatch.Stream)
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot initialize cloudwatch: %w", err)
	}

	return cloudWatchWriter, cloudWatchWriter.Close, nil
}

// InitializeStdout initializes logging to standard output with human-friendly output.
// It is used in unit and database tests.
func InitializeStdout() {
	configureZerolog()
	log.Logger = decorate(log.Output(stdoutWriter(true))).With().Str("binary", config.BinaryName()).Logger()
}

// InitializeLogger initializes logging to cloudwatch client and enables sentry Error logging.
// If cloudwatch is disabled, we enable stdout output.
func InitializeLogger() (zerolog.Logger, func()) {
	configureZerolog()

	var writers []io.Writer
	var closeFns []func()

	if config.Cloudwatch.Enabled {
		cwWriter, cwClose, err := cloudwatchWriter()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize cloudwatch")
			panic(err)
		}
		writers = append(writers, cwWriter)
		closeFns = append(closeFns, cwClose)
	} else {
		log.Trace().Msg("Cloudwatch not enabled, enabling stdout")
		writers = append(writers, stdoutWriter(false))
	}
	if config.Sentry.Dsn != "" {
		sWriter, closeFn, err := sentryWriter(config.Sentry.Dsn)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize sentry, disabling sentry monitoring")
		} else {
			writers = append(writers, sWriter)
			closeFns = append(closeFns, closeFn)
		}
	}

	// closes all writers in one func
	closeFn := func() {
		for _, fn := range closeFns {
			fn()
		}
	}
	return decorate(zerolog.New(io.MultiWriter(writers...))), closeFn
}

func DumpConfigForDevelopment() {
	// safe dump of app config
	if !config.InProdClowder() {
		config.DumpClowder(log.Logger)
		config.DumpConfig(log.Logger)
	}
}
