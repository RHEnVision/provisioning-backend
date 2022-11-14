package services

import (
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
		renderError(w, r, payloads.NewMissingRequestParameterError(r.Context(), "region parameter is missing"))
	}

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderNewErrorFromClientErr(w, r, err)
		return
	}

	authentication, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		renderNewErrorFromClientErr(w, r, err)
		return
	}

	ec2Client, err := clients.GetEC2Client(r.Context(), authentication, region)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS EC2 client", err))
		return
	}

	instances, err := ec2Client.ListInstanceTypesWithPaginator(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to list AWS EC2 instances", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListInstanceTypeResponse(instances)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render instance types list", err))
		return
	}
}
