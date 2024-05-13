package rbac

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/RHEnVision/provisioning-backend/internal/usrerr"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/rs/zerolog"
)

const TraceName = "github.com/EnVision/provisioning/internal/clients/http/rbac"

type rbac struct {
	client *ClientWithResponses
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
	ctx, span := telemetry.StartSpan(ctx, "Ready")
	defer span.End()

	logger := logger(ctx)
	resp, err := c.client.GetStatusWithResponse(ctx, headers.AddRbacIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Error().Err(err).Msgf("Readiness request failed for RBAC: %s", err.Error())
		return err
	}

	if resp == nil {
		return fmt.Errorf("ready call: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
		return fmt.Errorf("ready call: %w: %d", clients.ErrUnexpectedBackendResponse, resp.StatusCode())
	}

	return nil
}

// ErrMetaNotPresent is returned when metadata for pagination is not present in the response.
var ErrMetaNotPresent = fmt.Errorf("RBAC did not return metadata: %w", usrerr.ErrBadRequest400)

// FetchLimit is the maximum possible entries returned in one request
var FetchLimit = ptr.To(500)

func (c *rbac) GetPrincipalAccess(ctx context.Context) (clients.RbacAcl, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetPrincipalAccess")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	rhId := identity.Identity(ctx)
	orgID := rhId.Identity.OrgID
	accountNumber := rhId.Identity.AccountNumber
	var fullRequest bool
	var result clients.AccessList

	err := cache.Find(ctx, orgID+accountNumber, &result)
	if errors.Is(err, cache.ErrNotFound) {
		fullRequest = true
		start := time.Now()
		defer func() {
			metrics.RbacAclFetchDuration.Observe(float64(time.Since(start).Nanoseconds()) / 1000000)
		}()

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
			resp, gpaErr := c.client.GetPrincipalAccessWithResponse(ctx, &params, headers.AddRbacIdentityHeader, headers.AddEdgeRequestIdHeader)
			if gpaErr != nil {
				return nil, fmt.Errorf("get principal access: %w", gpaErr)
			}

			if resp == nil {
				return nil, fmt.Errorf("get principal access: empty response: %w", clients.ErrUnexpectedBackendResponse)
			}

			if resp.JSON200 == nil {
				return nil, fmt.Errorf("get principal access: %w", clients.ErrUnexpectedBackendResponse)
			}

			logger.Trace().
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
				logger.Warn().Msg("Maximum amount of RBAC requests reached, giving up")
				break
			}
		}

		err = cache.Set(ctx, orgID+accountNumber, &result)
		if err != nil {
			return nil, fmt.Errorf("acl cache set error: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("acl cache get error: %w", err)
	}

	logger.Debug().Msgf("ACL (cache: %t): %s", !fullRequest, result.String())
	return result, nil
}
