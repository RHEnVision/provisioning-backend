package services

import (
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var (
	UnknownProviderTypeError         = errors.New("unknown provider type parameter")
	ProviderTypeNotImplementedError  = errors.New("provider type not yet implemented")
	InvalidRequestPubkeyNewError     = errors.New("provide either existing (via NewName/NewBody) or new pubkey (ExistingID)")
	InvalidRequestPubkeyMissingError = errors.New("provide both NewName and NewBody for pubkey")
)

// CreateReservation dispatches requests to type provider specific handlers
func CreateReservation(w http.ResponseWriter, r *http.Request) {
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
		renderError(w, r, payloads.NewRenderError(r.Context(), "list reservations", err))
		return
	}
}

func GetReservationDetail(w http.ResponseWriter, r *http.Request) {
	id, err := ParseInt64(r, "ID")
	if err != nil {
		renderError(w, r, payloads.NewURLParsingError(r.Context(), "ID", err))
		return
	}

	rDao := dao.GetReservationDao(r.Context())

	// TODO: Add support for GCP and Azure, not generic reservation
	switch pType := models.ProviderTypeFromString(chi.URLParam(r, "TYPE")); pType {
	case models.ProviderTypeAWS:
		reservation, err := rDao.GetAWSById(r.Context(), id)
		if err != nil {
			renderNotFoundOrDAOError(w, r, err, "get reservation detail")
			return
		}

		if err := render.Render(w, r, payloads.NewAWSReservationResponse(reservation)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "reservation", err))
		}
	case models.ProviderTypeUnknown, models.ProviderTypeNoop:
		reservation, err := rDao.GetById(r.Context(), id)
		if err != nil {
			renderNotFoundOrDAOError(w, r, err, "get reservation detail")
			return
		}

		if err := render.Render(w, r, payloads.NewReservationResponse(reservation)); err != nil {
			renderError(w, r, payloads.NewRenderError(r.Context(), "reservation", err))
		}
	case models.ProviderTypeAzure:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	case models.ProviderTypeGCP:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	}
}
