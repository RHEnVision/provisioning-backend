package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/parsing"
)

type SourcesClient struct {
	client *ClientWithResponses
}

func init() {
	GetSourcesClientV2 = newSourcesClient
}

func newSourcesClient(ctx context.Context) (ClientV2, error) {
	c, err := NewClientWithResponses(config.Sources.URL)
	if err != nil {
		return nil, err
	}
	return &SourcesClient{client: c}, nil
}

func (c *SourcesClient) ListProvisioningSources(ctx context.Context) (*[]Source, error) {
	return nil, errors.New("not implemented")
}

func (c *SourcesClient) GetArn(ctx context.Context, sourceId string) (string, error) {
	return "", errors.New("not implemented")
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
	ctxval.Logger(ctx).Info().Msg("Fetching the Application Type ID of Provisioning for Sources")
	resp, err := c.client.ListApplicationTypes(ctx, &ListApplicationTypesParams{})
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return "", fmt.Errorf("failed to fetch ApplicationTypes: %w", err)
	}
	if !parsing.IsHTTPStatus2xx(resp.StatusCode) {
		ctxval.Logger(ctx).Warn().Msgf("Sources replied with unexpected status while fetching ApplicationTypes: %s", resp.Status)
		return "", fmt.Errorf("%w, status: '%s'", ApplicationTypesFetchUnsuccessful, resp.Status)
	}
	defer resp.Body.Close()
	var appTypesData dataElement
	if err = json.NewDecoder(resp.Body).Decode(&appTypesData); err != nil {
		return "", fmt.Errorf("could not unmarshal the response: %w", err)
	}
	for _, t := range appTypesData.Data {
		if t.Name == "/insights/platform/provisioning" {
			ctxval.Logger(ctx).Info().Msgf("The Application Type ID found: '%s' and it got cached", t.Id)
			return t.Id, nil
		}
	}
	return "", ApplicationTypeNotFound
}
