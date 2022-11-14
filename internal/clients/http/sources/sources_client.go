package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

const TraceName = "github.com/EnVision/provisioning/internal/clients/http/sources"

type sourcesClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetSourcesClient = newSourcesClient
}

func logger(ctx context.Context) zerolog.Logger {
	return ctxval.Logger(ctx).With().Str("client", "sources").Logger()
}

func newSourcesClient(ctx context.Context) (clients.Sources, error) {
	c, err := NewClientWithResponses(config.Sources.URL, func(c *Client) error {
		c.Client = http.NewPlatformClient(ctx, config.Sources.Proxy.URL)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &sourcesClient{client: c}, nil
}

type appType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type dataElement struct {
	Data []appType `json:"data"`
}

func (c *sourcesClient) Ready(ctx context.Context) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "Ready")
	defer span.End()

	logger := logger(ctx)
	resp, err := c.client.ListApplicationTypes(ctx, &ListApplicationTypesParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		logger.Error().Err(err).Msgf("Readiness request failed for sources: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		return fmt.Errorf("ready call: %w", err)
	}
	return nil
}

func (c *sourcesClient) ListProvisioningSources(ctx context.Context) ([]*clients.Source, error) {
	logger := logger(ctx)
	logger.Trace().Msg("Listing provisioning sources")

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get provisioning type id")
		return nil, fmt.Errorf("failed to get provisioning app type: %w", err)
	}

	resp, err := c.client.ListApplicationTypeSourcesWithResponse(ctx, appTypeId, &ListApplicationTypeSourcesParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return nil, fmt.Errorf("failed to get ApplicationTypes: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			return nil, fmt.Errorf("list provisioning sources call: %w", http.SourceNotFoundErr)
		}
		return nil, fmt.Errorf("list provisioning sources call: %w", err)
	}

	result := make([]*clients.Source, len(*resp.JSON200.Data))
	for i, src := range *resp.JSON200.Data {
		newSrc := clients.Source{
			Id:           src.Id,
			Name:         src.Name,
			SourceTypeId: src.SourceTypeId,
			Uid:          src.Uid,
		}
		result[i] = &newSrc
	}
	return result, nil
}

func (c *sourcesClient) GetAuthentication(ctx context.Context, sourceId clients.ID) (*clients.Authentication, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting authentication from source %s", sourceId)

	// Get all the authentications linked to a specific source
	resp, err := c.client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &ListSourceAuthenticationsParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		return nil, fmt.Errorf("cannot list source authentication: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			return nil, fmt.Errorf("get source authentication call: %w", http.AuthenticationForSourcesNotFoundErr)
		}
		return nil, fmt.Errorf("get source authentication call: %w", err)
	}

	// Filter authentications to include only auth where resource_type == "Application". We do this because
	// Sources API currently does not provide a good server-side filtering.
	auth, err := filterSourceAuthentications(resp.JSON200.Data)
	if err != nil {
		logger.Warn().Msgf("Sources replied with more than one authentications for source: %s", sourceId)
		return nil, err
	}

	// Get the resource_id which equals to application_id
	// and check that application_type_id in /applications/<app_id> equals to provisioning id
	res, err := c.client.ShowApplicationWithResponse(ctx, *auth.ResourceId, headers.AddSourcesIdentityHeader)
	if err != nil {
		return nil, fmt.Errorf("cannot list source authentication: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			return nil, fmt.Errorf("get source ARN call: %w", http.ApplicationNotFoundErr)
		}
		return nil, fmt.Errorf("get source ARN call: %w", err)
	}

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot get provisioning app type: %w", err)
	}

	if *res.JSON200.ApplicationTypeId != appTypeId {
		return nil, fmt.Errorf("%w for source id %s", http.AuthenticationSourceAssociationErr, sourceId)
	}

	authentication := clients.NewAuthenticationFromSourceAuthType(ctx, *auth.Username, string(*auth.Authtype))
	return authentication, nil
}

func (c *sourcesClient) GetProvisioningTypeId(ctx context.Context) (string, error) {
	appTypeId, err := cache.FindAppTypeId(ctx)
	if errors.Is(err, cache.NotFound) {
		appTypeId, err = c.loadAppId(ctx)
		if err != nil {
			return "", err
		}
		err = cache.SetAppTypeId(ctx, appTypeId)
		if err != nil {
			return "", fmt.Errorf("unable to store app type id to cache: %w", err)
		}
	} else if err != nil {
		return "", fmt.Errorf("unable to get app type id from cache: %w", err)
	}

	return appTypeId, nil
}

func (c *sourcesClient) loadAppId(ctx context.Context) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msg("Fetching the Application Type ID of Provisioning for Sources")

	resp, err := c.client.ListApplicationTypes(ctx, &ListApplicationTypesParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return "", fmt.Errorf("failed to fetch ApplicationTypes: %w", err)
	}
	defer resp.Body.Close()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			return "", fmt.Errorf("load app ID call: %w", http.ApplicationTypeNotFoundErr)
		}
		return "", fmt.Errorf("load app ID call: %w", err)
	}

	var appTypesData dataElement
	if err = json.NewDecoder(resp.Body).Decode(&appTypesData); err != nil {
		return "", fmt.Errorf("could not unmarshal application type response: %w", err)
	}
	for _, t := range appTypesData.Data {
		if t.Name == "/insights/platform/provisioning" {
			logger.Trace().Msgf("The Application Type ID found: '%s' and it got cached", t.Id)
			return t.Id, nil
		}
	}
	return "", http.ApplicationTypeNotFoundErr
}

func filterSourceAuthentications(authentications *[]AuthenticationRead) (AuthenticationRead, error) {
	auths := *authentications
	list := make([]AuthenticationRead, len(auths))
	for i, auth := range auths {
		if *auth.ResourceType == "Application" {
			list[i] = auth
		}
	}
	// Assumption: each source has one authentication linked to it
	if len(list) > 1 {
		return AuthenticationRead{}, http.MoreThanOneAuthenticationForSourceErr
	}
	return list[0], nil
}
