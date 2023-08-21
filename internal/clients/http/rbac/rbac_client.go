package rbac

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

const TraceName = "github.com/EnVision/provisioning/internal/clients/http/rbac"

type rbac struct {
	client *ClientWithResponses
}

func init() {
	clients.GetRbacClient = newClient
}

func logger(ctx context.Context) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("client", "ib").Logger()
}

func newClient(ctx context.Context) clients.Rbac {
	if !config.Application.RbacEnabled {
		return &allPermRbacSingleton
	}

	c, err := NewClientWithResponses(config.RBAC.URL, func(c *Client) error {
		c.Client = http.NewPlatformClient(ctx, config.ImageBuilder.Proxy.URL)
		return nil
	})
	if err != nil {
		// in case there was an error during initialization return "no permissions ACL"
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Could not initialize RBAC client, returning empty ACL")
		return &noPermRbacSingleton
	}
	return &rbac{client: c}
}

func (c *rbac) Ready(ctx context.Context) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "Ready")
	defer span.End()

	logger := logger(ctx)
	resp, err := c.client.GetStatus(ctx, headers.AddRbacIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Error().Err(err).Msgf("Readiness request failed for RBAC: %s", err.Error())
		return err
	}
	defer func() {
		if tempErr := resp.Body.Close(); tempErr != nil {
			logger.Error().Err(tempErr).Msg("Readiness request for RBAC: response body close error")
		}
	}()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		return fmt.Errorf("ready call: %w", err)
	}
	return nil
}

// ErrMetaNotPresent is returned when metadata for pagination is not present in the response.
var ErrMetaNotPresent = fmt.Errorf("RBAC did not return metadata: %w", clients.HttpClientErr)

// Maximum possible entries returned in one request
var FetchLimit = ptr.To(500)

func (c *rbac) GetPrincipalAccess(ctx context.Context) (clients.RbacAcl, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "GetPrincipalAccess")
	defer span.End()

	start := time.Now()
	defer func() {
		metrics.RbacAclFetchDuration.Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
	}()

	var result clients.AccessList
	records := math.MaxInt
	offset := 0
	maxQueries := 10

	// keep fetching until we have all the records
	for len(result) < records {
		params := GetPrincipalAccessParams{
			Application: "provisioning",
			Limit:       FetchLimit,
			Offset:      &offset,
		}
		resp, err := c.client.GetPrincipalAccessWithResponse(ctx, &params, headers.AddRbacIdentityHeader, headers.AddEdgeRequestIdHeader)
		if err != nil {
			return nil, fmt.Errorf("get principal access: %w", err)
		}

		err = http.HandleHTTPResponses(ctx, resp.StatusCode())
		if err != nil {
			return nil, fmt.Errorf("get principal access: %w", err)
		}
		zerolog.Ctx(ctx).Trace().
			Int("limit", *FetchLimit).
			Int("offset", offset).
			Int("length", len(result)).
			Int("entries", len(resp.JSON200.Data)).
			Int("return_code", resp.StatusCode()).
			Msg("Performed get principal access RBAC call")

		if resp.JSON200 == nil || resp.JSON200.Meta == nil || resp.JSON200.Meta.Count == nil {
			return nil, ErrMetaNotPresent
		}
		records = int(*resp.JSON200.Meta.Count)

		for _, a := range resp.JSON200.Data {
			result = append(result, clients.NewAccess(a.Permission))
		}
		offset += *FetchLimit

		maxQueries -= 1
		if maxQueries <= 0 {
			zerolog.Ctx(ctx).Warn().Msg("Maximum amount of RBAC requests reached, giving up")
			break
		}
	}

	zerolog.Ctx(ctx).Trace().Msgf("Access list: %s", result.String())
	return result, nil
}
