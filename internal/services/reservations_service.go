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

var UnknownProviderTypeError = errors.New("unknown provider type parameter")
var ProviderTypeNotImplementedError = errors.New("provider type not yet implemented")
var InvalidRequestPubkeyNewError = errors.New("provide either existing (via NewName/NewBody) or new pubkey (ExistingID)")
var InvalidRequestPubkeyMissingError = errors.New("provide both NewName and NewBody for pubkey")

// CreateReservation dispatches requests to type provider specific handlers
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	pType := models.ProviderTypeFromString(chi.URLParam(r, "type"))
	switch pType {
	case models.ProviderTypeNoop:
		CreateNoopReservation(w, r)
	case models.ProviderTypeAWS:
		CreateAWSReservation(w, r)
	case models.ProviderTypeAzure:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	case models.ProviderTypeGCE:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), ProviderTypeNotImplementedError))
	case models.ProviderTypeUnknown:
	default:
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), UnknownProviderTypeError))
	}
}

func ListReservations(w http.ResponseWriter, r *http.Request) {
	rDao, err := dao.GetReservationDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "reservation DAO", err))
		return
	}

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
