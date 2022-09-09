package services

import (
	"errors"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/lzap/dejq"
)

func CreateAWSReservation(w http.ResponseWriter, r *http.Request) {
	logger := ctxval.Logger(r.Context())

	var accountId int64 = ctxval.AccountId(r.Context())

	payload := &payloads.AWSReservationRequestPayload{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), err))
		return
	}

	rDao, err := dao.GetReservationDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "reservation DAO", err))
		return
	}
	pkDao, err := dao.GetPubkeyDao(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewInitializeDAOError(r.Context(), "pubkey DAO", err))
		return
	}

	detail := &models.AWSDetail{
		Region:       payload.Region,
		InstanceType: payload.InstanceType,
		Amount:       payload.Amount,
		PowerOff:     payload.PowerOff,
	}
	reservation := &models.AWSReservation{
		PubkeyID: payload.PubkeyID,
		SourceID: payload.SourceID,
		ImageID:  payload.ImageID,
		Detail:   detail,
	}
	reservation.AccountID = accountId
	reservation.Status = "Created"
	reservation.Provider = models.ProviderTypeAWS
	reservation.Steps = 2
	if payload.Name != nil {
		newName := config.AWS.InstancePrefix + *payload.Name
		reservation.Detail.Name = &newName
	}

	// validate pubkey
	logger.Debug().Msgf("Validating existence of pubkey %d for this account", reservation.PubkeyID)
	pk, err := pkDao.GetById(r.Context(), reservation.PubkeyID)
	if err != nil {
		var e dao.NoRowsError
		if errors.As(err, &e) {
			renderError(w, r, payloads.NewNotFoundError(r.Context(), err))
		} else {
			renderError(w, r, payloads.NewDAOError(r.Context(), "get pubkey by id", err))
		}
		return
	}
	logger.Debug().Msgf("Found pubkey %d named '%s'", pk.ID, pk.Name)

	// create reservation in the database
	err = rDao.CreateAWS(r.Context(), reservation)
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

	// Fetch arn from Sources
	arn, err := sourcesClient.GetArn(r.Context(), payload.SourceID)
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", err, 404))
			return
		}
		renderError(w, r, payloads.ClientError(r.Context(), "Sources", "can't fetch arn from sources", err, 500))
		return
	}

	// Prepare jobs
	logger.Debug().Msgf("Enqueuing upload key job for pubkey %d", pk.ID)
	uploadPubkeyJob := dejq.PendingJob{
		Type: queue.TypePubkeyUploadAws,
		Body: &jobs.PubkeyUploadAWSTaskArgs{
			AccountID:     accountId,
			ReservationID: reservation.ID,
			Region:        reservation.Detail.Region,
			PubkeyID:      pk.ID,
			ARN:           arn,
			SourceID:      reservation.SourceID,
		},
	}

	var ami string
	if strings.HasPrefix(reservation.ImageID, "ami-") {
		// Direct AMI ID was provided, no need to call image builder
		ami = reservation.ImageID
	} else {
		// Get Image builder client
		IBClient, ibErr := clients.GetImageBuilderClient(r.Context())
		logger.Trace().Msg("Creating IB client")
		if ibErr != nil {
			renderError(w, r, payloads.NewClientInitializationError(r.Context(), "image builder client", ibErr))
			return
		}

		// Get AMI
		ami, ibErr = IBClient.GetAWSAmi(r.Context(), reservation.ImageID)
		if ibErr != nil {
			renderError(w, r, payloads.ClientError(r.Context(), "Image Builder", "can't get ami from image builder", ibErr, 500))
		}
	}

	logger.Debug().Msgf("Enqueuing launch instance job for source %s", reservation.SourceID)
	launchJob := dejq.PendingJob{
		Type: queue.TypeLaunchInstanceAws,
		Body: &jobs.LaunchInstanceAWSTaskArgs{
			AccountID:     accountId,
			ReservationID: reservation.ID,
			Region:        reservation.Detail.Region,
			PubkeyID:      pk.ID,
			Detail:        reservation.Detail,
			AMI:           ami,
			ARN:           arn,
		},
	}

	startJobs := []dejq.PendingJob{uploadPubkeyJob, launchJob}

	// Enqueue all jobs
	err = queue.GetEnqueuer().Enqueue(r.Context(), startJobs...)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "AWS reservation", err))
		return
	}

	// Return response payload
	if err := render.Render(w, r, payloads.NewAWSReservationResponse(reservation)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "reservation", err))
	}
}
