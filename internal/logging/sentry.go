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
	sWriter, err := sentrywriter.New(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot initialize sentry: %w", err)
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("stream", config.BinaryName())
	})
	fWriter := NewSentryReplacer(sWriter)

	closeFunc := func() {
		_ = fWriter.Close()
		_ = sWriter.Close()
	}
	return fWriter, closeFunc, nil
}
