package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

func AzureOfferingTemplate(w http.ResponseWriter, r *http.Request) {
	clientName := config.Azure.ClientPrincipalName
	if clientName == "" {
		clientName = "Red Hat Launch images client"
	}

	tmpl := clients.AzureOfferingTemplate{
		OfferingDefaultName:        "Red Hat Hybrid Cloud Console",
		OfferingDefaultDescription: "Allows Red Hat to upload images and deploy Virtual Machines from Hybrid cloud console",
		TenantID:                   config.Azure.TenantID,
		PrincipalID:                config.Azure.ClientPrincipalID,
		PrincipalName:              clientName,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Allow", "GET, HEAD, OPTIONS")
	if err := tmpl.Render(r.Context(), w); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "failed to render the Azure template", err))
		return
	}
}
