package sources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/parsing"
)

type AppType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type dataElement struct {
	Data []AppType `json:"data"`
}

func getSourcesClient(ctx context.Context) (SourcesIntegration, error) {
	return NewClientWithResponses(config.Sources.URL)
}

func init() {
	GetSourcesClient = getSourcesClient
}

func (c *ClientWithResponses) GetProvisioningTypeId(ctx context.Context, reqEditors ...RequestEditorFn) (string, error) {
	if appTypeId, ok := cache.AppTypeId(); ok {
		return appTypeId, nil
	}
	appTypeId, err := loadAppId(ctx, c)
	if err != nil {
		return "", err
	}
	cache.SetAppTypeId(appTypeId)
	return appTypeId, nil
}

func loadAppId(ctx context.Context, client *ClientWithResponses) (string, error) {
	ctxval.Logger(ctx).Info().Msg("Fetching the Application Type ID of Provisioning for Sources")
	resp, err := client.ListApplicationTypes(ctx, &ListApplicationTypesParams{})
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Failed to fetch ApplicationTypes from sources")
		return "", fmt.Errorf("failed to fetch ApplicationTypes: %w", err)
	}
	if !parsing.IsHTTPStatus2xx(resp.StatusCode) {
		ctxval.Logger(ctx).Warn().Msgf("Sources replied with unexpected status while fetching ApplicationTypes: %s", resp.Status)
		return "", fmt.Errorf("%w, status: '%s'", ApplicationTypesFetchUnsuccessfulErr, resp.Status)
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
	return "", ApplicationTypeNotFoundErr
}
