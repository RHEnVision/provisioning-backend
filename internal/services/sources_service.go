package services

import (
	"net/http"

	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/RHEnVision/provisioning-backend/internal/parsing"
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
	appTypeId, err := client.GetProvisioningTypeId(r.Context(), headers.AddIdentityHeader)
	if err != nil {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "get provisioning app type", err, 500))
		return
	}
	resp, err := client.ListApplicationTypeSourcesWithResponse(r.Context(), appTypeId, &sources.ListApplicationTypeSourcesParams{}, headers.AddIdentityHeader)
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

	client, err := sources.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client", err))
		return
	}

	resp, err := client.ShowSourceWithResponse(r.Context(), sources.ID(id), headers.AddIdentityHeader)
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
