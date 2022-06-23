package services

import (
	"context"
	"fmt"
	"net/http"

	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	client, err := sources.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client", err))
		return
	}
	resp, err := client.ListApplicationTypeSourcesWithResponse(r.Context(), config.Sources.AppTypeId, &sources.ListApplicationTypeSourcesParams{}, AddIdentityHeader)
	if err != nil {
		renderError(w, r, payloads.New3rdPartyClientError(r.Context(), "list sources", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListSourcesResponse(resp.JSON200.Data)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list sources", err))
		return
	}

}

func GetSource(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ID")

	client, err := sources.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client", err))
		return
	}

	resp, err := client.ShowSourceWithResponse(r.Context(), sources.ID(id), AddIdentityHeader)
	if err != nil {
		renderError(w, r, payloads.New3rdPartyClientError(r.Context(), "show source", err))
		return
	}
	if err := render.Render(w, r, payloads.NewShowSourcesResponse(resp.JSON200)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "show source", err))
	}
}

func fetchARN(ctx context.Context, client sources.SourcesIntegration, sourceId string) (string, error) {
	// Get all the authentications linked to a specific source
	resp, err := client.ListSourceAuthenticationsWithResponse(ctx, sourceId, &sources.ListSourceAuthenticationsParams{}, AddIdentityHeader)
	if err != nil {
		return "", fmt.Errorf("cannot list source authentication: %w", err)
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

	if *res.JSON200.ApplicationTypeId == config.Sources.AppTypeId {
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
