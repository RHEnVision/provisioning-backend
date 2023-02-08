package http

import (
	"context"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Shared HTTP transport for all platform clients to utilize connection caching
var transport = &http.Transport{}

// NewPlatformClient returns new HTTP client (doer) with W3C Trace Context, logging tracing
// and/or HTTP proxy (non-clowder environment only) according to application configuration.
// Use this function to create HTTP clients for communication with all platform services.
func NewPlatformClient(ctx context.Context, proxy string) HttpRequestDoer {
	var rt http.RoundTripper = transport

	if proxy != "" {
		if config.InClowder() {
			ctxval.Logger(ctx).Warn().Msgf("Unable to use HTTP client proxy in clowder environment: %s", proxy)
		} else {
			ctxval.Logger(ctx).Warn().Msgf("Creating HTTP client with proxy %s", proxy)
			rt = &http.Transport{Proxy: http.ProxyURL(config.StringToURL(proxy))}
		}
	}

	if config.Telemetry.Enabled {
		rt = otelhttp.NewTransport(rt)
	}

	var doer HttpRequestDoer = &http.Client{Transport: rt}
	if config.RestEndpoints.TraceData {
		doer = NewLoggingDoer(ctx, doer)
	}
	return doer
}
