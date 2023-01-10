package services

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/lzap/dejq"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

func CreateGCPReservation(w http.ResponseWriter, r *http.Request) {
	logger := *ctxval.Logger(r.Context())

	var accountId int64 = ctxval.AccountId(r.Context())

	payload := &payloads.GCPReservationRequestPayload{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "GCP reservation", err))
		return
	}

	rDao := dao.GetReservationDao(r.Context())
	pkDao := dao.GetPubkeyDao(r.Context())

	detail := &models.GCPDetail{
		Zone:        payload.Zone,
		MachineType: payload.MachineType,
		Amount:      payload.Amount,
		PowerOff:    payload.PowerOff,
	}
	reservation := &models.GCPReservation{
		PubkeyID: payload.PubkeyID,
		ImageID:  payload.ImageID,
		SourceID: payload.SourceID,
		Detail:   detail,
	}

	reservation.AccountID = accountId
	reservation.Status = "Created"
	reservation.Provider = models.ProviderTypeGCP
	reservation.Steps = 1

	logger.Debug().Msgf("Validating existence of pubkey %d for this account", reservation.PubkeyID)
	pk, err := pkDao.GetById(r.Context(), reservation.PubkeyID)
	if err != nil {
		message := fmt.Sprintf("get pubkey with id %d", reservation.PubkeyID)
		renderNotFoundOrDAOError(w, r, err, message)
		return
	}
	logger.Debug().Msgf("Found pubkey %d named '%s'", pk.ID, pk.Name)

	// create reservation in the database
	err = rDao.CreateGCP(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create reservation", err))
		return
	}
	logger = logger.With().Int64("reservation_id", reservation.ID).Logger()
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	// Get Sources client
	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	// Fetch project id from Sources
	authentication, err := sourcesClient.GetAuthentication(r.Context(), payload.SourceID)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	if typeErr := authentication.MustBe(models.ProviderTypeGCP); typeErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), typeErr))
		return
	}

	// TODO: upload key job if needed

	// Get Image builder client
	IBClient, ibErr := clients.GetImageBuilderClient(r.Context())
	logger.Trace().Msg("Creating IB client")
	if ibErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), ibErr))
		return
	}

	// Get Image Name
	name, ibErr := IBClient.GetGCPImageName(r.Context(), reservation.ImageID)
	if ibErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), ibErr))
		return
	}

	logger.Trace().Msgf("Image Name is %s", name)

	launchJob := dejq.PendingJob{
		Type: queue.TypeLaunchInstanceGcp,
		Body: &jobs.LaunchInstanceGCPTaskArgs{
			AccountID:     accountId,
			ReservationID: reservation.ID,
			Zone:          reservation.Detail.Zone,
			PubkeyID:      reservation.PubkeyID,
			Detail:        reservation.Detail,
			ImageName:     name,
			ProjectID:     authentication,
		},
	}
	logger.Debug().Interface("job", launchJob).Msgf("Enqueuing launch instance job for source %s", reservation.SourceID)

	startJobs := []dejq.PendingJob{launchJob}

	// Enqueue all jobs
	err = queue.GetEnqueuer().Enqueue(r.Context(), startJobs...)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "enqueing task GCP reservation error", err))
		return
	}

	// Return response payload
	if err := render.Render(w, r, payloads.NewGCPReservationResponse(reservation)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		return
	}
}
