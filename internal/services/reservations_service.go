package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	UnknownProviderTypeError        = errors.New("unknown provider type parameter")
	ProviderTypeMismatchError       = errors.New("reservation type does not match requested provider type")
	ProviderTypeNotImplementedError = errors.New("provider type not yet implemented")
	UnknownInstanceTypeNameError    = errors.New("unknown instance type")
	ArchitectureMismatch            = errors.New("instance type and image architecture mismatch")
	BothTypeAndTemplateMissingError = errors.New("instance type or launch template not set")
	UnsupportedRegionError          = errors.New("unknown region/location/zone")
)

// CreateReservation dispatches requests to type provider specific handlers
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	if !config.LaunchEnabled(r.Context()) {
		writeUnauthorized(w, r)
	}

	pType := models.ProviderTypeFromString(chi.URLParam(r, "TYPE"))
	switch pType {
	case models.ProviderTypeNoop:
		CreateNoopReservation(w, r)
	case models.ProviderTypeAWS:
		CreateAWSReservation(w, r)
	case models.ProviderTypeAzure:
		if config.FeatureEnabled(r.Context(), "azure") {
			CreateAzureReservation(w, r)
		} else {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "azure reservation is not implemented", ProviderTypeNotImplementedError))
		}
	case models.ProviderTypeGCP:
		CreateGCPReservation(w, r)
	case models.ProviderTypeUnknown:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", UnknownProviderTypeError))
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", UnknownProviderTypeError))
	}
}

func ListReservations(w http.ResponseWriter, r *http.Request) {
	rDao := dao.GetReservationDao(r.Context())

	reservations, err := rDao.List(r.Context(), 100, 0)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "list reservations", err))
		return
	}

	if err := render.RenderList(w, r, payloads.NewReservationListResponse(reservations)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservations list", err))
		return
	}
}

func GetReservationDetail(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "TYPE")
	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "unable to parse ID parameter", err))
		return
	}

	// Get generic reservation and find its type
	rDao := dao.GetReservationDao(r.Context())
	reservation, err := rDao.GetById(r.Context(), id)
	if err != nil {
		renderNotFoundOrDAOError(w, r, err, "get reservation detail")
		return
	}

	providerType := models.ProviderTypeFromString(provider)
	if providerType != models.ProviderTypeUnknown && reservation.Provider != providerType {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider type", ProviderTypeMismatchError))
		return
	}

	switch providerType {
	// Generic reservation request will have provider == "" and thus render this
	case models.ProviderTypeUnknown, models.ProviderTypeNoop:
		if err := render.Render(w, r, payloads.NewReservationResponse(reservation)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	case models.ProviderTypeAWS:
		reservationAws, err := rDao.GetAWSById(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get AWS reservation with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		instances, err := rDao.ListInstances(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get reservation with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		if err := render.Render(w, r, payloads.NewAWSReservationResponse(reservationAws, instances)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	case models.ProviderTypeAzure:
		reservationAzure, err := rDao.GetAzureById(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get Azure reservation with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		instances, err := rDao.ListInstances(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get reservation instances with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		if err := render.Render(w, r, payloads.NewAzureReservationResponse(reservationAzure, instances)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	case models.ProviderTypeGCP:
		reservationGCP, err := rDao.GetGCPById(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get GCP reservation with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		instances, err := rDao.ListInstances(r.Context(), id)
		if err != nil {
			message := fmt.Sprintf("get reservation with id %d", id)
			renderNotFoundOrDAOError(w, r, err, message)
			return
		}

		if err := render.Render(w, r, payloads.NewGCPReservationResponse(reservationGCP, instances)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "provider is not supported", ProviderTypeNotImplementedError))
	}
}
