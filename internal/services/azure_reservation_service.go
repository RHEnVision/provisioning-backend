package services

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/logging"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
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
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func CreateAzureReservation(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	payload := &payloads.AzureReservationRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Azure reservation", err))
		return
	}

	pkDao := dao.GetPubkeyDao(r.Context())
	rDao := dao.GetReservationDao(r.Context())

	// validate region
	// it needs to be in region format, but also
	if payload.Location != "" && !preload.AzureInstanceType.ValidateRegion(payload.Location+"_1") {
		// Location format is accepted now for backwards compatibility, but we should deprecate it
		if preload.AzureInstanceType.ValidateRegion(payload.Location) {
			logger.Warn().Msgf("Azure region passed with location suffix (%s), this is deprecated behaviour format", payload.Location)
		} else {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Unsupported location", ErrUnsupportedRegion))
			return
		}
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
	resourceGroupName := payload.ResourceGroup
	// Azure image IDs are "free form", if it's a UUID we treat it like a compose ID
	if composeUUID, pErr := uuid.Parse(payload.ImageID); pErr == nil {
		// Composer-built image
		instanceType := preload.AzureInstanceType.FindInstanceType(clients.InstanceTypeName(payload.InstanceSize))
		if instanceType == nil {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Instance size is not a valid Azure instance size", nil))
			return
		}

		var imageResourceGroupName string
		imageResourceGroupName, azureImageName, err = ibClient.GetAzureImageInfo(r.Context(), composeUUID, *instanceType)
		if err != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), err))
			return
		}
		if resourceGroupName == "" {
			resourceGroupName = imageResourceGroupName
		}
		azureImageName = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/images/%s", authentication.Payload, imageResourceGroupName, azureImageName)
	} else {
		// Format Image ID for image names passed manually in here.
		// Assumes the image is in the resource group we want to deploy into.
		if strings.HasPrefix(payload.ImageID, "composer-api") {
			if resourceGroupName == "" {
				resourceGroupName = jobs.DefaultAzureResourceGroupName
			}
			azureImageName = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/images/%s", authentication.Payload, resourceGroupName, payload.ImageID)
		} else {
			// Anything else is treated like a direct Azure image ID (e.g. from https://imagedirectory.cloud)
			azureImageName = payload.ImageID
		}
	}

	supportedArch := "x86_64"
	it := preload.AzureInstanceType.FindInstanceType(clients.InstanceTypeName(payload.InstanceSize))
	if it == nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), fmt.Sprintf("unknown instance size: %s", payload.InstanceSize), ErrUnknownInstanceTypeName))
		return
	}
	if it.Architecture.String() != supportedArch {
		renderError(w, r, payloads.NewWrongArchitectureUserError(r.Context(), ErrArchitectureMismatch))
		return
	}

	name := config.Application.InstancePrefix + payload.Name
	detail := &models.AzureDetail{
		Location:      payload.Location,
		ResourceGroup: resourceGroupName,
		InstanceSize:  payload.InstanceSize,
		Amount:        payload.Amount,
		PowerOff:      payload.PowerOff,
		Name:          name,
	}
	reservation := &models.AzureReservation{
		PubkeyID: &payload.PubkeyID,
		SourceID: payload.SourceID,
		ImageID:  payload.ImageID,
		Detail:   detail,
	}
	reservation.Steps = int32(len(jobs.LaunchInstanceAzureSteps))
	reservation.StepTitles = jobs.LaunchInstanceAzureSteps

	// The last step: create reservation in the database and submit new job
	err = rDao.CreateAzure(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create Azure reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	launchJob := worker.Job{
		Type:      jobs.TypeLaunchInstanceAzure,
		Identity:  identity.Identity(r.Context()),
		EdgeID:    logging.EdgeRequestId(r.Context()),
		AccountID: identity.AccountId(r.Context()),
		Args: jobs.LaunchInstanceAzureTaskArgs{
			Location:          reservation.Detail.Location,
			ReservationID:     reservation.ID,
			ResourceGroupName: reservation.Detail.ResourceGroup,
			PubkeyID:          pk.ID,
			SourceID:          reservation.SourceID,
			AzureImageID:      azureImageName,
			Subscription:      authentication,
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
	if err = render.Render(w, r, payloads.NewAzureReservationResponse(reservation, unused)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render Azure reservation", err))
	}
}
