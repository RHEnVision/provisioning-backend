package services

import (
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
	resp, err := client.ListApplicationTypeSourcesWithResponse(r.Context(), config.Sources.AppId, &sources.ListApplicationTypeSourcesParams{}, AddIdentityHeader)
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
