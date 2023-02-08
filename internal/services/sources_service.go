package services

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

func ListSources(w http.ResponseWriter, r *http.Request) {
	var sourcesList []*clients.Source

	client, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}
	provider := r.URL.Query().Get("provider")
	asProviderType := models.ProviderTypeFromString(provider)
	if asProviderType != models.ProviderTypeUnknown {
		sourcesList, err = client.ListProvisioningSourcesByProvider(r.Context(), asProviderType)
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
	} else if provider == "" {
		sourcesList, err = client.ListAllProvisioningSources(r.Context())
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
	} else {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), fmt.Sprintf("unknown provider: %s", provider), clients.UnknownProviderErr))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListSourcesResponse(sourcesList)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render sources list", err))
		return
	}
}
