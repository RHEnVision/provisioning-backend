package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
)

// CreateNoopReservation is used to create empty reservation that is processed without any operation
// being made. This is useful when testing the job queue. The endpoint has no payload.
func CreateNoopReservation(w http.ResponseWriter, r *http.Request) {
	logger := ctxval.Logger(r.Context())

	// TODO: get this from X-RH-Identity via middleware/context
	var accountId int64 = 1

	rDao, err := dao.GetReservationDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "reservation DAO", err))
		return
	}
	reservation := &models.NoopReservation{
		Reservation: models.Reservation{
			Provider:  models.ProviderTypeNoop,
			AccountID: accountId,
			Status:    "Created",
		},
	}

	// create reservation in the database
	err = rDao.CreateNoop(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	// create a new job
	args := &jobs.NoopJobArgs{
		AccountID:     accountId,
		ReservationID: reservation.ID,
	}
	err = jobs.EnqueueNoop(r.Context(), args)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "EnqueueNoop", err))
		return
	}

	if err := render.Render(w, r, payloads.NewNoopReservationResponse(reservation)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "reservation", err))
	}
}
