package services

import (
	"net/http"

	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	client, err := sources.GetSourcesClientV2(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client", err))
		return
	}
	sourcesList, err := client.ListProvisioningSources(r.Context())
	if err != nil {
		renderError(w, r, payloads.SourcesClientError(r.Context(), "sources client error", err, 500))
		return
	}
	if err := render.RenderList(w, r, payloads.NewListSourcesResponse(sourcesList)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list sources", err))
		return
	}

}
