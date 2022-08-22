package services

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/lzap/dejq"
)

// CreateNoopReservation is used to create empty reservation that is processed without any operation
// being made. This is useful when testing the job queue. The endpoint has no payload.
func CreateNoopReservation(w http.ResponseWriter, r *http.Request) {
	logger := ctxval.Logger(r.Context())
	accountId := ctxval.AccountId(r.Context())

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
	pj := dejq.PendingJob{
		Type: queue.TypeNoop,
		Body: &jobs.NoopJobArgs{
			AccountID:     accountId,
			ReservationID: reservation.ID,
		},
	}
	logger.Debug().Interface("arg", pj.Body).Msgf("Enqueuing no operation job: %+v", pj.Body)
	err = queue.GetEnqueuer().Enqueue(r.Context(), pj)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "EnqueueNoop", err))
		return
	}

	if err := render.Render(w, r, payloads.NewNoopReservationResponse(reservation)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "reservation", err))
	}
}
