package logging

import (
	"fmt"
	"io"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	sentrywriter "github.com/archdx/zerolog-sentry"
	"github.com/getsentry/sentry-go"
)

// sentryWriter creates a zerolog writer for sentry.
// Uses github.com/archdx/zerolog-sentry which is very simple wrapper.
func sentryWriter(dsn string) (io.Writer, func(), error) {
	var opts []sentrywriter.WriterOption
	if config.Sentry.Debug {
		opts = append(opts, sentrywriter.WithDebug())
	}
	opts = append(opts, sentrywriter.WithSampleRate(config.Sentry.SampleRate))
	// TODO this calls sentry.Init again (it is initialized in our telemetry package)
	// We must be able to initialize it only once and pass it into here, or, create it
	// here with tracing and sample rate enabled and use it also for otel.
	// Currently it is not possible to have sentry for both errors and tracing.
	wr, err := sentrywriter.New(dsn, opts...)
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot initialize sentry: %w", err)
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("stream", config.BinaryName())
	})
	// close function returns error, but it's always nil, so we are ignoring it here
	return wr, func() { _ = wr.Close() }, nil
}
