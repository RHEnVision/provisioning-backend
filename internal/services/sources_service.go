package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/parsing"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client", err))
		return
	}
	appTypeId, err := client.GetProvisioningTypeId(r.Context(), AddIdentityHeader)
	if err != nil {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "get provisioning app type", err, 500))
		return
	}
	resp, err := client.ListApplicationTypeSourcesWithResponse(r.Context(), appTypeId, &sources.ListApplicationTypeSourcesParams{}, AddIdentityHeader)
	statusCode := resp.StatusCode()
	if err != nil {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "list sources", err, statusCode))
		return
	}
	if parsing.IsHTTPNotFound(statusCode) {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "sources not found", err, statusCode))
		return
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "sources client error", err, statusCode))
		return
	}
	if err := render.RenderList(w, r, payloads.NewListSourcesResponse(resp.JSON200.Data)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list sources", err))
		return
	}

}

func GetSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client", err))
		return
	}

	resp, err := client.ShowSourceWithResponse(r.Context(), sources.ID(id), AddIdentityHeader)
	statusCode := resp.StatusCode()
	if err != nil {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "show source", err, statusCode))
		return
	}
	if parsing.IsHTTPNotFound(statusCode) {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "source not found", err, statusCode))
		return
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "sources client error", err, statusCode))
		return
	}
	if err := render.Render(w, r, payloads.NewShowSourcesResponse(resp.JSON200)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "show source", err))
	}
}

func fetchARN(ctx context.Context, client clients.SourcesIntegration, sourceId string) (string, error) {
	// Get all the authentications linked to a specific source
	resp, err := client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &sources.ListSourceAuthenticationsParams{}, AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
	}
	statusCode := resp.StatusCode()
	if parsing.IsHTTPNotFound(statusCode) {
		return "", sources.AuthenticationForSourcesNotFoundErr
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		return "", sources.SourcesClientErr
	}
	// Filter authentications to include only auth where resource_type == "Application"
	auth, err := filterSourceAuthentications(resp.JSON200.Data)
	if err != nil {
		return "", err
	}
	// Get the resource_id which equals to application_id
	// and check that application_type_id in /applications/<app_id> equals to provisioning id
	res, err := client.ShowApplicationWithResponse(ctx, *auth.ResourceId, AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
	}
	statusCode = res.StatusCode()
	if parsing.IsHTTPNotFound(statusCode) {
		return "", sources.ApplicationNotFoundErr
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		return "", sources.SourcesClientErr
	}

	appTypeId, err := client.GetProvisioningTypeId(ctx, AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot get provisioning app type: %w", err)
	}

	if *res.JSON200.ApplicationTypeId == appTypeId {
		return *auth.Username, nil

	}
	return "", fmt.Errorf("cannot find authentication linked to source id %s and to the provisioning app: %w", sourceId, err)
}

func filterSourceAuthentications(authentications *[]sources.AuthenticationRead) (sources.AuthenticationRead, error) {
	auths := *authentications
	list := make([]sources.AuthenticationRead, 0, len(auths))
	for _, auth := range auths {
		if *auth.ResourceType == "Application" {
			list = append(list, auth)
		}
	}
	// Assumption: each source has one authentication linked to it
	if len(list) > 1 {
		return sources.AuthenticationRead{}, sources.MoreThenOneAuthenticationForSourceErr
	}
	return list[0], nil
}
