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
	"github.com/RHEnVision/provisioning-backend/internal/models"
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
	return NewSourcesClientWithUrl(ctx, config.Sources.URL)
}

// NewSourcesClientWithUrl allows customization of the URL for the underlying client.
// It is meant for testing only, for production please use clients.GetSourcesClient.
func NewSourcesClientWithUrl(ctx context.Context, url string) (clients.Sources, error) {
	c, err := NewClientWithResponses(url, func(c *Client) error {
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
	resp, err := c.client.ListApplicationTypes(ctx, &ListApplicationTypesParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Error().Err(err).Msg("Readiness request failed for sources")
		return err
	}
	defer resp.Body.Close()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		return fmt.Errorf("ready call: %w", err)
	}
	return nil
}

func (c *sourcesClient) ListProvisioningSourcesByProvider(ctx context.Context, provider models.ProviderType) ([]*clients.Source, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Listing provisioning sources of provider %s", provider)

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get provisioning type id")
		return nil, fmt.Errorf("failed to get provisioning app type: %w", err)
	}

	resp, err := c.client.ListApplicationTypeSourcesWithResponse(ctx, appTypeId, &ListApplicationTypeSourcesParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
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

	result := make([]*clients.Source, 0, len(*resp.JSON200.Data))

	for _, src := range *resp.JSON200.Data {
		sourceTypeName, err := c.GetSourceTypeName(ctx, *src.SourceTypeId)
		if err != nil {
			return nil, fmt.Errorf("could not get source type name for source type id %d: %w", src.SourceTypeId, err)
		}

		if sourceTypeName == models.ProviderTypeFromString(provider.String()) {
			newSrc := clients.Source{
				Id:           src.Id,
				Name:         src.Name,
				SourceTypeId: src.SourceTypeId,
				Uid:          src.Uid,
			}
			result = append(result, &newSrc)
		}
	}

	return result, nil
}

func (c *sourcesClient) ListAllProvisioningSources(ctx context.Context) ([]*clients.Source, error) {
	logger := logger(ctx)
	logger.Trace().Msg("Listing all provisioning sources")

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get provisioning type id")
		return nil, fmt.Errorf("failed to get provisioning app type: %w", err)
	}

	resp, err := c.client.ListApplicationTypeSourcesWithResponse(ctx, appTypeId, &ListApplicationTypeSourcesParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
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
	resp, err := c.client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &ListSourceAuthenticationsParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
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
	auth, err := filterSourceAuthentications(*resp.JSON200.Data)
	if err != nil {
		at := make([]string, 0)
		for _, auth := range *resp.JSON200.Data {
			at = append(at, string(*auth.Authtype))
		}
		logger.Warn().Msgf("Sources did not return any Provisioning authentication for source(auth types): %s(%v)", sourceId, at)
		return nil, err
	}

	authentication := clients.NewAuthenticationFromSourceAuthType(ctx, *auth.Username, string(*auth.Authtype), *auth.ResourceId)
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

	resp, err := c.client.ListApplicationTypes(ctx, &ListApplicationTypesParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
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

//nolint:exhaustive
func filterSourceAuthentications(authentications []AuthenticationRead) (AuthenticationRead, error) {
	for _, auth := range authentications {
		if *auth.ResourceType == "Application" {
			switch *auth.Authtype {
			// Type of the authentication as stored in Sources by listing the source types or the application types
			case "provisioning-arn",
				"provisioning_lighthouse_subscription_id",
				"provisioning_project_id":
				return auth, nil
			default:
				continue
			}
		}
	}
	return AuthenticationRead{}, http.ApplicationReadErr
}

func (c *sourcesClient) GetSourceTypeName(ctx context.Context, sourceTypeID string) (models.ProviderType, error) {
	logger := logger(ctx)
	logger.Trace().Msg("Getting source types list from sources")

	// Get all the source types
	resp, err := c.client.ListSourceTypesWithResponse(ctx, &ListSourceTypesParams{}, headers.AddSourcesIdentityHeader)
	if err != nil {
		return models.ProviderTypeUnknown, fmt.Errorf("cannot list source types: %w", err)
	}

	for _, st := range *resp.JSON200.Data {
		if sourceTypeID == *st.Id {
			logger.Trace().Msg("Found source type id from sources")
			switch *st.Name {
			case "amazon":
				return models.ProviderTypeAWS, nil
			case "google":
				return models.ProviderTypeGCP, nil
			case "azure":
				return models.ProviderTypeAzure, nil
			default:
				return models.ProviderTypeUnknown, fmt.Errorf("provider unknown %w", clients.UnknownProviderErr)
			}
		}
	}
	return models.ProviderTypeUnknown, fmt.Errorf("cannot find source type name for source type id %s: %w", sourceTypeID, http.SourceTypeNameNotFoundErr)
}
