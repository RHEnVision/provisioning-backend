package services

import (
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
	logger := ctxval.Logger(r.Context())

	var accountId int64 = ctxval.AccountId(r.Context())

	payload := &payloads.GCPReservationRequestPayload{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
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
		renderNotFoundOrDAOError(w, r, err, "get pubkey by id", "")
		return
	}
	logger.Debug().Msgf("Found pubkey %d named '%s'", pk.ID, pk.Name)

	// create reservation in the database
	err = rDao.CreateGCP(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	// Get Sources client
	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "sources client v2", err))
		return
	}

	// Fetch project id from Sources
	authentication, err := sourcesClient.GetAuthentication(r.Context(), payload.SourceID)
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch project id from sources", err, 404))
			return
		}
		renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch project id from sources", err, 500))
		return
	}

	if typeErr := authentication.MustBe(models.ProviderTypeGCP); typeErr != nil {
		renderError(w, r, payloads.ClientError(r.Context(), "Sources", "unexpected source type", typeErr, 500))
		return
	}

	// TODO: upload key job if needed

	// Get Image builder client
	IBClient, ibErr := clients.GetImageBuilderClient(r.Context())
	logger.Trace().Msg("Creating IB client")
	if ibErr != nil {
		renderError(w, r, payloads.NewClientInitializationError(r.Context(), "image builder client", ibErr))
		return
	}

	// Get Image Name
	name, ibErr := IBClient.GetGCPImageName(r.Context(), reservation.ImageID)
	if ibErr != nil {
		renderError(w, r, payloads.ClientError(r.Context(), "Image Builder", "can't get name from image builder", ibErr, 500))
		return
	}

	logger.Trace().Msgf("Image Name is %s", name)

	logger.Debug().Msgf("Enqueuing launch instance job for source %s", reservation.SourceID)
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

	startJobs := []dejq.PendingJob{launchJob}

	// Enqueue all jobs
	err = queue.GetEnqueuer().Enqueue(r.Context(), startJobs...)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "GCP reservation", err))
		return
	}

	// Return response payload
	if err := render.Render(w, r, payloads.NewGCPReservationResponse(reservation)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "reservation", err))
		return
	}
}
