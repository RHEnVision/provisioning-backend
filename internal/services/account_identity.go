package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetAccountIdentity(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	authentication, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	ec2Client, err := clients.GetEC2Client(r.Context(), authentication, "")
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS EC2 client", err))
		return
	}

	accountId, err := ec2Client.GetAccountId(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get account id", err))
		return
	}

	if err := render.Render(w, r, payloads.NewAccountIdentityResponse(accountId)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render account id", err))
		return
	}
}
