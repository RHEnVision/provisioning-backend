package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"net/url"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

const TraceName = telemetry.TracePrefix + "internal/clients/http/sources"

type sourcesClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetSourcesClient = newSourcesClient
}

func logger(ctx context.Context) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("client", "sources").Logger()
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
	resp, err := c.client.ListApplicationTypesWithResponse(ctx, &ListApplicationTypesParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Error().Err(err).Msg("Readiness request failed for sources")
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

func (c *sourcesClient) ListProvisioningSourcesByProvider(ctx context.Context, provider models.ProviderType) ([]*clients.Source, int, error) {
	logger := logger(ctx)
	params := &ListApplicationTypeSourcesParams{}
	ctx, span := otel.Tracer(TraceName).Start(ctx, "ListProvisioningSourcesByProvider")
	defer span.End()

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get provisioning type id")
		return nil, 0, fmt.Errorf("failed to get provisioning app type: %w", err)
	}

	sourcesProviderName := provider.SourcesProviderName()
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get provider name according to sources service")
		return nil, 0, fmt.Errorf("failed to get provider name according to sources service: %w", err)
	}

	offset := page.Offset(ctx).String()
	limit := page.Limit(ctx).String()

	resp, err := c.client.ListApplicationTypeSourcesWithResponse(ctx, appTypeId, params, headers.AddSourcesIdentityHeader,
		headers.AddEdgeRequestIdHeader, BuildQuery("filter[source_type][name]", sourcesProviderName, "offset", offset, "limit", limit))
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return nil, 0, fmt.Errorf("failed to get ApplicationTypes: %w", err)
	}

	if resp == nil {
		return nil, 0, fmt.Errorf("failed to get ApplicationTypes: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200 == nil {
		return nil, 0, fmt.Errorf("failed to get ApplicationTypes: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200.Data == nil {
		return nil, 0, fmt.Errorf("list provisioning sources call: %w", clients.ErrNoResponseData)
	}

	result := make([]*clients.Source, 0, len(*resp.JSON200.Data))

	for _, src := range *resp.JSON200.Data {
		newSrc := clients.Source{
			ID:           ptr.From(src.Id),
			Name:         ptr.From(src.Name),
			SourceTypeID: ptr.From(src.SourceTypeId),
			Uid:          ptr.From(src.Uid),
			Status:       string(*src.AvailabilityStatus),
		}
		result = append(result, &newSrc)
	}

	total := 0
	if resp.JSON200.Meta != nil {
		total = *resp.JSON200.Meta.Count
	}

	return result, total, nil
}

func (c *sourcesClient) ListAllProvisioningSources(ctx context.Context) ([]*clients.Source, int, error) {
	logger := logger(ctx)
	params := &ListApplicationTypeSourcesParams{}
	ctx, span := otel.Tracer(TraceName).Start(ctx, "ListAllProvisioningSources")
	defer span.End()

	appTypeId, err := c.GetProvisioningTypeId(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to get provisioning type id")
		return nil, 0, fmt.Errorf("failed to get provisioning app type: %w", err)
	}

	params.Offset = page.Offset(ctx).IntPtr()
	params.Limit = page.Limit(ctx).IntPtr()

	resp, err := c.client.ListApplicationTypeSourcesWithResponse(ctx, appTypeId, params, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return nil, 0, fmt.Errorf("failed to get ApplicationTypes: %w", err)
	}
	if resp == nil {
		return nil, 0, fmt.Errorf("list provisioning sources call: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200 == nil {
		return nil, 0, fmt.Errorf("list provisioning sources call: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200.Data == nil {
		return nil, 0, fmt.Errorf("list provisioning sources call: %w", clients.ErrNoResponseData)
	}

	result := make([]*clients.Source, len(*resp.JSON200.Data))

	for i, src := range *resp.JSON200.Data {
		newSrc := clients.Source{
			ID:           ptr.From(src.Id),
			Name:         ptr.From(src.Name),
			SourceTypeID: ptr.From(src.SourceTypeId),
			Uid:          ptr.From(src.Uid),
			Status:       string(*src.AvailabilityStatus),
		}
		result[i] = &newSrc
	}

	total := 0
	if resp.JSON200.Meta != nil {
		total = *resp.JSON200.Meta.Count
	}

	return result, total, nil
}

func (c *sourcesClient) GetAuthentication(ctx context.Context, sourceId string) (*clients.Authentication, error) {
	logger := logger(ctx)
	ctx, span := otel.Tracer(TraceName).Start(ctx, "GetAuthentication")
	defer span.End()

	// Get all the authentications linked to a specific source
	resp, err := c.client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &ListSourceAuthenticationsParams{}, headers.AddSourcesIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		return nil, fmt.Errorf("cannot list source authentication: %w", err)
	}

	// Filter authentications to include only auth where resource_type == "Application". We do this because
	// Sources API currently does not provide a good server-side filtering.

	if resp == nil {
		return nil, fmt.Errorf("get source authentication call: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get source authentication call: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200.Data == nil {
		return nil, fmt.Errorf("get source authentication call: %w", clients.ErrNoResponseData)
	}
	auth, err := filterSourceAuthentications(*resp.JSON200.Data)
	if err != nil {
		at := zerolog.Arr()
		for _, auth := range *resp.JSON200.Data {
			if auth.Authtype != nil {
				at.Str(*auth.Authtype)
			}
		}
		// Likely a super-key that hasn't been processed by super-key-worker yet
		logger.Warn().Str("source_id", sourceId).RawJSON("response", resp.Body).Array("filtered", at).Msg("Source does not have Provisioning authentication")
		return nil, err
	}

	if auth.Username == nil || auth.Authtype == nil || auth.ResourceId == nil {
		return nil, fmt.Errorf("cannot create source from source authentication type: %w", http.ErrSourcesInvalidAuthentication)
	}
	authentication, err := clients.NewAuthenticationFromSourceAuthType(ctx, *auth.Username, *auth.Authtype, *auth.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("cannot create source from source authentication type: %w", err)
	}
	return authentication, nil
}

func (c *sourcesClient) GetProvisioningTypeId(ctx context.Context) (string, error) {
	appTypeId, err := cache.FindAppTypeId(ctx)
	if errors.Is(err, cache.ErrNotFound) {
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

	if resp == nil {
		return "", fmt.Errorf("failed to fetch ApplicationTypes: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	defer func() {
		if tempErr := resp.Body.Close(); tempErr != nil {
			logger.Error().Err(tempErr).Msg("ApplicationTypes fetching from sources: response body close error")
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("failed to fetch ApplicationTypes: %w: %d", clients.ErrUnexpectedBackendResponse, resp.StatusCode)
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
	return "", http.ErrApplicationTypeNotFound
}

func BuildQuery(keysAndValues ...string) func(ctx context.Context, req *stdhttp.Request) error {
	return func(ctx context.Context, req *stdhttp.Request) error {
		if len(keysAndValues)%2 != 0 {
			panic("cannot build sources query: invalid input")
		}
		queryParams := make([]string, 0)
		for i := 0; i < len(keysAndValues); i += 2 {
			key := url.QueryEscape(keysAndValues[i])
			value := url.QueryEscape(keysAndValues[i+1])
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, value))
		}

		req.URL.RawQuery = strings.Join(queryParams, "&")
		return nil
	}
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
	return AuthenticationRead{}, http.ErrApplicationRead
}
