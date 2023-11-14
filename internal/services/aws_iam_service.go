package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/payloads/validation"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func ValidatePermissions(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())
	sourceId := chi.URLParam(r, "ID")
	if err := validation.DigitsOnly(sourceId); err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "id parameter invalid", err))
	}

	region := r.URL.Query().Get("region")

	if region == "" {
		region = config.AWS.DefaultRegion
	}

	// Get Sources client
	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	// Fetch arn from Sources
	authentication, err := sourcesClient.GetAuthentication(r.Context(), sourceId)
	if err != nil {
		if err != nil {
			if errors.Is(err, clients.ErrNotFound) {
				renderError(w, r, payloads.NewNotFoundError(r.Context(), "unable to get authentication for sources", err))
				return
			}
			if errors.Is(err, clients.ErrBadRequest) {
				renderError(w, r, payloads.NewResponseError(r.Context(), http.StatusBadRequest, "unable to get authentication from sources", err))
				return
			}
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	if !authentication.Is(models.ProviderTypeAWS) {
		if err = render.Render(w, r, payloads.NewPermissionsResponse(nil)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render missing permissions", err))
			return
		}
		return
	}

	ec2Client, err := clients.GetEC2Client(r.Context(), authentication, region)
	if err != nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS EC2 client", err))
		return
	}

	logger.Info().Msgf("Listing permissions.")
	missingPermissions, err := ec2Client.CheckPermission(r.Context(), authentication)
	if err != nil && missingPermissions == nil {
		renderError(w, r, payloads.NewAWSError(r.Context(), "unable to check aws permissions", err))
		return
	}

	if err := render.Render(w, r, payloads.NewPermissionsResponse(missingPermissions)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render missing permissions", err))
		return
	}
}
