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
	accountId := ctxval.AccountId(r.Context())

	rDao := dao.GetReservationDao(r.Context())
	reservation := &models.NoopReservation{
		Reservation: models.Reservation{
			Provider:   models.ProviderTypeNoop,
			AccountID:  accountId,
			Status:     "Created",
			Steps:      1,
			StepTitles: []string{"A test step"},
		},
	}

	// create reservation in the database
	err := rDao.CreateNoop(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create noop reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	// create a new job
	args := jobs.NoopJobArgs{
		AccountID:     accountId,
		ReservationID: reservation.ID,
	}
	err = jobs.EnqueueNoop(r.Context(), args)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "enqueuing noop reservation error", err))
		return
	}

	if err := render.Render(w, r, payloads.NewNoopReservationResponse(reservation)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
	}
}
