package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func ListInstanceTypes(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")
	region := r.URL.Query().Get("region")
	if region == "" {
		renderError(w, r, payloads.NewMissingRequestParameterError(r.Context(), "region"))
	}

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "can't init sources client", err))
		return
	}

	authentication, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", err, 404))
			return
		}
		renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", err, 500))
		return
	}

	ec2Client, err := clients.GetCustomerEC2Client(r.Context(), authentication, region)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "failed to establish ec2 connection", err))
		return
	}

	instances, err := ec2Client.ListInstanceTypesWithPaginator(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "can't list EC2 instance types", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListInstanceTypeResponse(instances)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "list instance types", err))
		return
	}
}
