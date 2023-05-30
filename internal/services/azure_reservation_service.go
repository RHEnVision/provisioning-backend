package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/RHEnVision/provisioning-backend/internal/preload"
	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func CreateAzureReservation(w http.ResponseWriter, r *http.Request) {
	logger := *zerolog.Ctx(r.Context())

	payload := &payloads.AzureReservationRequestPayload{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Azure reservation", err))
		return
	}

	pkDao := dao.GetPubkeyDao(r.Context())
	rDao := dao.GetReservationDao(r.Context())

	// Check for preloaded region
	if payload.Location == "" {
		payload.Location = "eastus_1"
	}
	if !preload.AzureInstanceType.ValidateRegion(payload.Location) {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Unsupported location", UnsupportedRegionError))
		return
	}

	// Validate pubkey
	logger.Debug().Msgf("Validating existence of pubkey %d for this account", payload.PubkeyID)
	pk, err := pkDao.GetById(r.Context(), payload.PubkeyID)
	if err != nil {
		message := fmt.Sprintf("get pubkey with id %d", payload.PubkeyID)
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

	// Get IB client
	ibClient, err := clients.GetImageBuilderClient(r.Context())
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	// Fetch SubscriptionID from Sources
	authentication, err := sourcesClient.GetAuthentication(r.Context(), payload.SourceID)
	if err != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), err))
		return
	}

	if typeErr := authentication.MustBe(models.ProviderTypeAzure); typeErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), typeErr))
		return
	}

	var azureImageName string
	// Azure image IDs are "free form", if it's a UUID we treat it like a compose ID
	if _, pErr := uuid.Parse(payload.ImageID); pErr == nil {
		// Composer-built image
		azureImageName, err = ibClient.GetAzureImageID(r.Context(), payload.ImageID)
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
		azureImageName = fmt.Sprintf("/subscriptions/%s%s", authentication.Payload, azureImageName)
	} else {
		// Format Image ID for image names passed manually in here.
		// Assumes 'redhat-deployed' resource group.
		if strings.HasPrefix(payload.ImageID, "composer-api") {
			azureImageName = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/images/%s", authentication.Payload, "redhat-deployed", payload.ImageID)
		} else {
			// Anything else is treated like a direct Azure image ID (e.g. from https://imagedirectory.cloud)
			azureImageName = payload.ImageID
		}
	}

	supportedArch := "x86_64"
	it := preload.AzureInstanceType.FindInstanceType(clients.InstanceTypeName(payload.InstanceSize))
	if it == nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), fmt.Sprintf("unknown instance size: %s", payload.InstanceSize), UnknownInstanceTypeNameError))
		return
	}
	if it.Architecture.String() != supportedArch {
		renderError(w, r, payloads.NewWrongArchitectureUserError(r.Context(), ArchitectureMismatch))
		return
	}

	name := config.Application.InstancePrefix + payload.Name
	detail := &models.AzureDetail{
		Location:     payload.Location,
		InstanceSize: payload.InstanceSize,
		Amount:       payload.Amount,
		PowerOff:     payload.PowerOff,
		Name:         name,
	}
	reservation := &models.AzureReservation{
		PubkeyID: payload.PubkeyID,
		SourceID: payload.SourceID,
		ImageID:  payload.ImageID,
		Detail:   detail,
	}
	reservation.Steps = int32(len(jobs.LaunchInstanceAzureSteps))
	reservation.StepTitles = jobs.LaunchInstanceAzureSteps

	// create reservation in the database
	err = rDao.CreateAzure(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create Azure reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	launchJob := worker.Job{
		Type:      jobs.TypeLaunchInstanceAzure,
		Identity:  ctxval.Identity(r.Context()),
		AccountID: ctxval.AccountId(r.Context()),
		Args: jobs.LaunchInstanceAzureTaskArgs{
			ReservationID: reservation.ID,
			Location:      reservation.Detail.Location,
			PubkeyID:      pk.ID,
			SourceID:      reservation.SourceID,
			AzureImageID:  azureImageName,
			Subscription:  authentication,
		},
	}

	err = queue.GetEnqueuer(r.Context()).Enqueue(r.Context(), &launchJob)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "job enqueue error", err))
		return
	}

	// Return response payload
	unused := make([]*models.ReservationInstance, 0, 0)
	if err = render.Render(w, r, payloads.NewAzureReservationResponse(reservation, unused)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render Azure reservation", err))
	}
}
