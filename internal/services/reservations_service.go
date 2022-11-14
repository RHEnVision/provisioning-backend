package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	UnknownProviderTypeError         = errors.New("unknown provider type parameter")
	ProviderTypeMismatchError        = errors.New("reservation type does not match requested provider type")
	ProviderTypeNotImplementedError  = errors.New("provider type not yet implemented")
	InvalidRequestPubkeyNewError     = errors.New("provide either existing (via NewName/NewBody) or new pubkey (ExistingID)")
	InvalidRequestPubkeyMissingError = errors.New("provide both NewName and NewBody for pubkey")
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
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	case models.ProviderTypeGCP:
		CreateGCPReservation(w, r)
	case models.ProviderTypeUnknown:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), UnknownProviderTypeError))
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), UnknownProviderTypeError))
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
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeMismatchError))
		return
	}

	switch providerType {
	// Generic reservation request will have provider == "" and thus render this
	case models.ProviderTypeUnknown, models.ProviderTypeNoop:
		if err := render.Render(w, r, payloads.NewReservationResponse(reservation)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	case models.ProviderTypeAWS:
		reservation, err := rDao.GetAWSById(r.Context(), id)
		if err != nil {
			renderNotFoundOrDAOError(w, r, err, "get reservation detail")
			return
		}

		instances, err := rDao.ListInstances(r.Context(), id)
		if err != nil {
			renderNotFoundOrDAOError(w, r, err, "get reservation detail")
			return
		}

		if err := render.Render(w, r, payloads.NewAWSReservationResponse(reservation, instances)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		}
	case models.ProviderTypeAzure:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	case models.ProviderTypeGCP:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	}
}
