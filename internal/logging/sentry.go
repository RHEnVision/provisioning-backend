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
	wr, err := sentrywriter.New(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("cannot initialize sentry: %w", err)
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("stream", config.BinaryName())
	})
	// close function returns error, but it's always nil, so we are ignoring it here
	return wr, func() { _ = wr.Close() }, nil
}
