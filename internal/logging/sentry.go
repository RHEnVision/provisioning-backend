package logging

import (
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	sentrywriter "github.com/archdx/zerolog-sentry"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

// sentryWriter creates a zerolog writer for sentry.
// Uses github.com/archdx/zerolog-sentry which is very simple wrapper.
func sentryWriter(dsn string) (zerolog.LevelWriter, func(), error) {
	replacer := NewSentryReplacer()

	sWriter, err := sentrywriter.New(dsn, sentrywriter.WithBeforeSend(replacer.Replace))
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot initialize sentry: %w", err)
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("stream", config.BinaryName())
	})

	return sWriter, func() { _ = sWriter.Close }, nil
}
