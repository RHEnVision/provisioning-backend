package telemetry

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ClientProxyProductionUseErr is returned when HTTP proxy is configured in a Clowder environment
var ClientProxyProductionUseErr = errors.New("HTTP client proxy cannot be used in clowder")

// HTTPClient returns newly initialized http.Client type with OpenTelemetry
// transport. This client should be used for all REST communication on the
// platform.
//
// An optional HTTP proxy can be passed in, nil means no proxy.
func HTTPClient(ctx context.Context, proxy *url.URL) (*http.Client, error) {
	logger := ctxval.Logger(ctx)
	var proxyFunc func(*http.Request) (*url.URL, error)
	if proxy != nil {
		if config.InClowder() {
			return nil, ClientProxyProductionUseErr
		}
		logger.Trace().Msgf("Creating HTTP client with proxy: %s", proxy.String())
		proxyFunc = http.ProxyURL(proxy)
	} else {
		logger.Trace().Msg("Creating HTTP client")
	}

	// TODO: make timeouts configurable (taken from the std library)
	var transport http.RoundTripper
	transport = &http.Transport{
		Proxy: proxyFunc,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if config.Telemetry.Enabled {
		transport = otelhttp.NewTransport(transport)
	}
	c := http.Client{Transport: transport}
	return &c, nil
}
