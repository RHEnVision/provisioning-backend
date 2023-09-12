package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/RHEnVision/provisioning-backend/internal/preload"
	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

func CreateAWSReservation(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	var accountId int64 = identity.AccountId(r.Context())
	var id identity.Principal = identity.Identity(r.Context())

	payload := &payloads.AWSReservationRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "AWS reservation", err))
		return
	}

	rDao := dao.GetReservationDao(r.Context())
	pkDao := dao.GetPubkeyDao(r.Context())

	// Check for preloaded region
	if payload.Region == "" {
		payload.Region = "us-east-1"
	}
	if !preload.EC2InstanceType.ValidateRegion(payload.Region) {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Unsupported region", UnsupportedRegionError))
		return
	}

	// Either Launch Template or Instance Type must be set. Both can be set too, in that case, instance type overrides the launch template.
	if payload.InstanceType == "" && payload.LaunchTemplateID == "" {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Both instance type and launch template are missing", BothTypeAndTemplateMissingError))
		return
	}

	// Validate architecture match (hardcoded since image builder currently only supports x86_64). This can be only done
	// when launch template is not set.
	if payload.LaunchTemplateID == "" {
		supportedArch := "x86_64"
		it := preload.EC2InstanceType.FindInstanceType(clients.InstanceTypeName(payload.InstanceType))
		if it == nil {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), fmt.Sprintf("unknown type: %s", payload.InstanceType), UnknownInstanceTypeNameError))
			return
		}
		if it.Architecture.String() != supportedArch {
			renderError(w, r, payloads.NewWrongArchitectureUserError(r.Context(), ArchitectureMismatch))
			return
		}
	}

	detail := &models.AWSDetail{
		Region:           payload.Region,
		LaunchTemplateID: payload.LaunchTemplateID,
		InstanceType:     payload.InstanceType,
		Amount:           payload.Amount,
		PowerOff:         payload.PowerOff,
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
	reservation.Steps = 3
	reservation.StepTitles = []string{"Ensure public key", "Launch instance(s)", "Fetch instance(s) description"}
	newName := config.Application.InstancePrefix + payload.Name
	reservation.Detail.Name = &newName

	// validate pubkey - must be always present because of data integrity (foreign keys)
	logger.Debug().Msgf("Validating existence of pubkey %d for this account", reservation.PubkeyID)
	pk, err := pkDao.GetById(r.Context(), reservation.PubkeyID)
	if err != nil {
		message := fmt.Sprintf("get pubkey with id %d", reservation.PubkeyID)
		renderNotFoundOrDAOError(w, r, err, message)
		return
	}
	logger.Debug().Msgf("Found pubkey %d named '%s'", pk.ID, pk.Name)

	// Get Sources client
	sourcesClient, err := clients.GetSourcesClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	// Fetch arn from Sources
	authentication, err := sourcesClient.GetAuthentication(r.Context(), payload.SourceID)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	if typeErr := authentication.MustBe(models.ProviderTypeAWS); typeErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), typeErr))
		return
	}

	var ami string
	if reservation.ImageID == "" || strings.HasPrefix(reservation.ImageID, "ami-") {
		// Direct AMI or no image were provided (launch template), no need to call image builder
		ami = reservation.ImageID
	} else {
		// Not prefixed with "ami-" therefore this must be a valid UUID
		// Get Image builder client
		IBClient, ibErr := clients.GetImageBuilderClient(r.Context())
		logger.Trace().Msg("Creating IB client")
		if ibErr != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), ibErr))
			return
		}

		// Get AMI
		ami, ibErr = IBClient.GetAWSAmi(r.Context(), reservation.ImageID)
		if ibErr != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), ibErr))
			return
		}
	}

	// The last step: create reservation in the database and submit new job
	err = rDao.CreateAWS(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	launchJob := worker.Job{
		Type:      jobs.TypeLaunchInstanceAws,
		Identity:  id,
		AccountID: accountId,
		Args: jobs.LaunchInstanceAWSTaskArgs{
			ReservationID:    reservation.ID,
			Region:           reservation.Detail.Region,
			PubkeyID:         pk.ID,
			SourceID:         reservation.SourceID,
			Detail:           reservation.Detail,
			AMI:              ami,
			LaunchTemplateID: reservation.Detail.LaunchTemplateID,
			ARN:              authentication,
		},
	}

	err = queue.GetEnqueuer(r.Context()).Enqueue(r.Context(), &launchJob)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "job enqueue error", err))
		return
	}
	logger.Debug().Msgf("Enqueued reservation job %s", launchJob.ID)

	// Return response payload
	unused := make([]*models.ReservationInstance, 0, 0)
	if err := render.Render(w, r, payloads.NewAWSReservationResponse(reservation, unused)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render AWS reservation", err))
	}
}
