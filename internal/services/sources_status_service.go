package services

import (
	"errors"
	stdhttp "net/http"

	"github.com/RHEnVision/provisioning-backend/internal/payloads/validation"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
)

var ErrUnknownProviderFromSources = errors.New("unknown provider returned from sources")

// SourcesStatus fetches information from sources and then performs a smallest possible
// request on the cloud provider (list keys or similar). Reports an error if sources configuration
// is no longer valid.
func SourcesStatus(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	sourceId := chi.URLParam(r, "ID")
	if err := validation.DigitsOnly(sourceId); err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "id parameter invalid", err))
	}

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

	var statusClient clients.ClientStatuser
	switch auth.Type() {
	case models.ProviderTypeAWS:
		statusClient, err = clients.GetEC2Client(r.Context(), auth, config.AWS.DefaultRegion)
		if err != nil {
			renderError(w, r, payloads.NewAWSError(r.Context(), "unable to get AWS client", err))
			return
		}
	case models.ProviderTypeGCP:
		statusClient, err = clients.GetGCPClient(r.Context(), auth)
		if err != nil {
			renderError(w, r, payloads.NewGCPError(r.Context(), "unable to get GCP client", err))
			return
		}
	case models.ProviderTypeAzure:
		statusClient, err = clients.GetAzureClient(r.Context(), auth)
		if err != nil {
			renderError(w, r, payloads.NewAzureError(r.Context(), "unable to get Azure client", err))
			return
		}
	case models.ProviderTypeNoop:
	case models.ProviderTypeUnknown:
	default:
		renderError(w, r, payloads.NewStatusError(r.Context(), "unknown sources provider", ErrUnknownProviderFromSources))
		return
	}

	err = statusClient.Status(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewStatusError(r.Context(), "client status error", err))
		return
	}

	writeOk(w, r)
}
