package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

//nolint:exhaustive
func ListLaunchTemplates(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")

	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	auth, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	switch auth.ProviderType {
	case models.ProviderTypeAWS:
		ListLaunchTemplateAWS(w, r)
	case models.ProviderTypeAzure:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "azure reservation is not implemented", ProviderTypeNotImplementedError))
	case models.ProviderTypeGCP:
		ListLaunchTemplateGCP(w, r)
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", UnknownProviderTypeError))
	}
}

func ListLaunchTemplateAWS(w http.ResponseWriter, r *http.Request) {
	sourceId := chi.URLParam(r, "ID")
	region := r.URL.Query().Get("region")
	if region == "" {
		renderError(w, r, payloads.NewMissingRequestParameterError(r.Context(), "region parameter is missing"))
	}

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

	ec2Client, err := clients.GetEC2Client(r.Context(), authentication, region)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS EC2 client", err))
		return
	}

	templates, err := ec2Client.ListLaunchTemplates(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to list AWS EC2 launch templates", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListLaunchTemplateResponse(templates)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render launch templates list", err))
		return
	}
}

func ListLaunchTemplateGCP(w http.ResponseWriter, r *http.Request) {
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

	gcpClient, err := clients.GetGCPClient(r.Context(), authentication)
	if err != nil {
		renderError(w, r, payloads.NewGCPError(r.Context(), "unable to get GCP client", err))
		return
	}

	templates, err := gcpClient.ListLaunchTemplates(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewGCPError(r.Context(), "unable to list GCP launch templates", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewListLaunchTemplateResponse(templates)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render launch templates list", err))
		return
	}
}
