package sources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/headers"

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

func (client *ClientWithResponses) FetchARN(ctx context.Context, sourceId string) (string, error) {
	// Get all the authentications linked to a specific source
	resp, err := client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &ListSourceAuthenticationsParams{}, headers.AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
	}
	statusCode := resp.StatusCode()
	if parsing.IsHTTPNotFound(statusCode) {
		return "", AuthenticationForSourcesNotFoundErr
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		return "", SourcesClientErr
	}
	// Filter authentications to include only auth where resource_type == "Application"
	auth, err := client.FilterSourceAuthentications(resp.JSON200.Data)
	if err != nil {
		return "", err
	}
	// Get the resource_id which equals to application_id
	// and check that application_type_id in /applications/<app_id> equals to provisioning id
	res, err := client.ShowApplicationWithResponse(ctx, *auth.ResourceId, headers.AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
	}
	statusCode = res.StatusCode()
	if parsing.IsHTTPNotFound(statusCode) {
		return "", ApplicationNotFoundErr
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		return "", SourcesClientErr
	}

	appTypeId, err := client.GetProvisioningTypeId(ctx, headers.AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot get provisioning app type: %w", err)
	}

	if *res.JSON200.ApplicationTypeId == appTypeId {
		return *auth.Username, nil

	}
	return "", fmt.Errorf("cannot find authentication linked to source id %s and to the provisioning app: %w", sourceId, err)
}

func (client *ClientWithResponses) FilterSourceAuthentications(authentications *[]AuthenticationRead) (AuthenticationRead, error) {
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
