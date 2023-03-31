package telemetry

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/rs/zerolog"
)

// LoggingTransport logs warning when HTTP reply is not 2xx/3xx. In addition,
// when trace logging level is enabled, it also dumps request/reply body.
type LoggingTransport struct {
	roundTripper http.RoundTripper
	logger       *zerolog.Logger
}

func (d *LoggingTransport) RoundTrip(h *http.Request) (*http.Response, error) {
	var dump []byte
	if d.logger.GetLevel() <= zerolog.TraceLevel {
		dump, _ = httputil.DumpRequestOut(h, true)
	}
	resp, err := d.roundTripper.RoundTrip(h)
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		d.logger.Warn().Msgf("Sentry replied with error: %s", resp.Status)
		if d.logger.GetLevel() <= zerolog.TraceLevel {
			d.logger.Trace().Msgf("Sentry HTTP request: %s", dump)
			dump, _ = httputil.DumpResponse(resp, true)
			d.logger.Trace().Msgf("Sentry HTTP response: %s", dump)
		}
	}
	if err != nil {
		err = fmt.Errorf("roundtrip error: %w", err)
	}
	return resp, err
}
