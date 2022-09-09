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
)

// TODO This should have been not exported
type SourcesClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetSourcesClient = newSourcesClient
}

func newSourcesClient(ctx context.Context) (clients.Sources, error) {
	c, err := NewClientWithResponses(config.Sources.URL, func(c *Client) error {
		if config.Sources.Proxy.URL != "" {
			if config.Features.Environment != "development" {
				return clients.ClientProxyProductionUseErr
			}
			var client HttpRequestDoer
			client, err := clients.NewProxyDoer(ctx, config.Sources.Proxy.URL)
			if err != nil {
				return fmt.Errorf("cannot create proxy doer: %w", err)
			}
			if config.RestEndpoints.TraceData {
				client = clients.NewLoggingDoer(ctx, client)
			}
			c.Client = client
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &SourcesClient{client: c}, nil
}

func copySource(src Source) clients.Source {
	return clients.Source{
		Id:           src.Id,
		Name:         src.Name,
		SourceTypeId: src.SourceTypeId,
		Uid:          src.Uid,
	}
}

type appType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type dataElement struct {
	Data []appType `json:"data"`
}

func (c *SourcesClient) Ready(ctx context.Context) error {
	logger := ctxval.Logger(ctx)

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

func (c *SourcesClient) ListProvisioningSources(ctx context.Context) (*[]clients.Source, error) {
	logger := ctxval.Logger(ctx)
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
		if errors.Is(err, clients.NotFoundError) {
			return nil, SourceNotFoundErr
		}
		return nil, fmt.Errorf("list provisioning sources call: %w", err)
	}

	result := make([]clients.Source, 0, len(*resp.JSON200.Data))
	for _, s := range *resp.JSON200.Data {
		result = append(result, copySource(s))
	}
	return &result, nil
}

func (c *SourcesClient) GetArn(ctx context.Context, sourceId clients.ID) (string, error) {
	logger := ctxval.Logger(ctx)
	logger.Trace().Msgf("Getting ARN of source %v", sourceId)

	// Get all the authentications linked to a specific source
	resp, err := c.client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &ListSourceAuthenticationsParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundError) {
			return "", SourceNotFoundErr
		}
		return "", fmt.Errorf("get source ARN call: %w", err)
	}

	// Filter authentications to include only auth where resource_type == "Application". We do this because
	// Sources API currently does not provide a good server-side filtering.
	auth, err := filterSourceAuthentications(resp.JSON200.Data)
	if err != nil {
		logger.Warn().Msgf("Sources replied with more then one authenticatios for source: %vs", sourceId)
		return "", err
	}

	// Get the resource_id which equals to application_id
	// and check that application_type_id in /applications/<app_id> equals to provisioning id
	res, err := c.client.ShowApplicationWithResponse(ctx, *auth.ResourceId, headers.AddSourcesIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundError) {
			return "", ApplicationNotFoundErr
		}
		return "", fmt.Errorf("get source ARN call: %w", err)
	}

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		return "", fmt.Errorf("cannot get provisioning app type: %w", err)
	}

	if *res.JSON200.ApplicationTypeId != appTypeId {
		return "", fmt.Errorf("%w for source id %s", AuthenticationSourceAssociationErr, sourceId)
	}
	return *auth.Username, nil
}

func (c *SourcesClient) GetProvisioningTypeId(ctx context.Context) (string, error) {
	if appTypeId, ok := cache.AppTypeId(); ok {
		return appTypeId, nil
	}
	appTypeId, err := c.loadAppId(ctx)
	if err != nil {
		return "", err
	}
	cache.SetAppTypeId(appTypeId)
	return appTypeId, nil
}

func (c *SourcesClient) loadAppId(ctx context.Context) (string, error) {
	logger := ctxval.Logger(ctx)
	logger.Trace().Msg("Fetching the Application Type ID of Provisioning for Sources")

	resp, err := c.client.ListApplicationTypes(ctx, &ListApplicationTypesParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return "", fmt.Errorf("failed to fetch ApplicationTypes: %w", err)
	}
	defer resp.Body.Close()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		if errors.Is(err, clients.NotFoundError) {
			return "", ApplicationTypeNotFoundErr
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
	return "", ApplicationTypeNotFoundErr
}

func filterSourceAuthentications(authentications *[]AuthenticationRead) (AuthenticationRead, error) {
	auths := *authentications
	list := make([]AuthenticationRead, 0, len(auths))
	for _, auth := range auths {
		if *auth.ResourceType == "Application" {
			list = append(list, auth)
		}
	}
	// Assumption: each source has one authentication linked to it
	if len(list) > 1 {
		return AuthenticationRead{}, MoreThenOneAuthenticationForSourceErr
	}
	return list[0], nil
}
