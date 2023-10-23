package services

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/logging"

	"github.com/RHEnVision/provisioning-backend/internal/preload"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/RHEnVision/provisioning-backend/internal/queue"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/go-chi/render"
)

func CreateGCPReservation(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	var accountId int64 = identity.AccountId(r.Context())
	var id identity.Principal = identity.Identity(r.Context())

	payload := &payloads.GCPReservationRequest{}
	if err := render.Bind(r, payload); err != nil {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "GCP reservation", err))
		return
	}

	rDao := dao.GetReservationDao(r.Context())
	pkDao := dao.GetPubkeyDao(r.Context())

	// Check for preloaded region
	if !preload.GCPInstanceType.ValidateRegion(payload.Zone) {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Unsupported zone", ErrUnsupportedRegion))
		return
	}

	// Check machine type for "This organization policy prevents creating instances with exotic machine types. Contact the IT Public Cloud team at help.redhat.com for an exception"
	if payload.MachineType == "" {
		renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Machine type must be present even when a template is in use", ErrBothTypeAndTemplateMissing))
		return
	}

	namePattern := "inst-####"
	// Verify name pattern is lower cased and add #####
	if payload.NamePattern != "" {
		ok := isValidNamePattern(payload.NamePattern)
		if ok {
			namePattern = fmt.Sprintf("%s-#####", payload.NamePattern)
		} else {
			renderError(w, r, payloads.NewInvalidRequestError(r.Context(), "Invalid name pattern", ErrInvalidNamePattern))
			return
		}
	}

	resUUID := uuid.New().String()
	detail := &models.GCPDetail{
		NamePattern:      &namePattern,
		Zone:             payload.Zone,
		MachineType:      payload.MachineType,
		Amount:           payload.Amount,
		PowerOff:         payload.PowerOff,
		UUID:             resUUID,
		LaunchTemplateID: payload.LaunchTemplateID,
	}
	reservation := &models.GCPReservation{
		PubkeyID: &payload.PubkeyID,
		ImageID:  payload.ImageID,
		SourceID: payload.SourceID,
		Detail:   detail,
	}

	reservation.AccountID = accountId
	reservation.Status = "Created"
	reservation.Provider = models.ProviderTypeGCP
	reservation.Steps = 2
	reservation.StepTitles = jobs.LaunchInstanceGCPSteps

	if reservation.PubkeyID == nil {
		renderError(w, r, payloads.NewNotFoundError(r.Context(), "could not create AWS reservation", ErrPubkeyNotFound))
	}

	logger.Debug().Msgf("Validating existence of pubkey %d for this account", *reservation.PubkeyID)
	pk, err := pkDao.GetById(r.Context(), *reservation.PubkeyID)
	if err != nil {
		message := fmt.Sprintf("get pubkey with id %d", *reservation.PubkeyID)
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

	// Get Image builder client
	ibc, ibErr := clients.GetImageBuilderClient(r.Context())
	logger.Trace().Msg("Creating IB client")
	if ibErr != nil {
		renderError(w, r, payloads.NewClientError(r.Context(), ibErr))
		return
	}

	// Validate image
	var name string
	if _, pErr := uuid.Parse(payload.ImageID); pErr == nil {
		// Composer-built image
		name, ibErr = ibc.GetGCPImageName(r.Context(), reservation.ImageID)
		if ibErr != nil {
			renderError(w, r, payloads.NewClientError(r.Context(), ibErr))
			return
		}

		logger.Trace().Msgf("Image Name is %s", name)
	} else {
		// Treat HTTP(S) URLs like direct image ID (e.g. from https://imagedirectory.cloud)
		name = payload.ImageID
	}

	// The last step: create reservation in the database and submit new job
	err = rDao.CreateGCP(r.Context(), reservation)
	if err != nil {
		renderError(w, r, payloads.NewDAOError(r.Context(), "create reservation", err))
		return
	}
	logger.Debug().Msgf("Created a new reservation %d", reservation.ID)

	launchJob := worker.Job{
		Type:      jobs.TypeLaunchInstanceGcp,
		AccountID: accountId,
		EdgeID:    logging.EdgeRequestId(r.Context()),
		Identity:  id,
		Args: jobs.LaunchInstanceGCPTaskArgs{
			ReservationID:    reservation.ID,
			Zone:             reservation.Detail.Zone,
			PubkeyID:         *reservation.PubkeyID,
			Detail:           reservation.Detail,
			ImageName:        name,
			ProjectID:        authentication,
			LaunchTemplateID: reservation.Detail.LaunchTemplateID,
		},
	}

	err = queue.GetEnqueuer(r.Context()).Enqueue(r.Context(), &launchJob)
	if err != nil {
		renderError(w, r, payloads.NewEnqueueTaskError(r.Context(), "job enqueue error", err))
		return
	}
	logger.Debug().Msgf("Enqueued reservation job %s", launchJob.ID)

	unused := make([]*models.ReservationInstance, 0, 0)
	// Return response payload
	if err := render.Render(w, r, payloads.NewGCPReservationResponse(reservation, unused)); err != nil {
		renderError(w, r, payloads.NewRenderError(r.Context(), "unable to render reservation", err))
		return
	}
}

func isValidNamePattern(namePattern string) bool {
	if strings.ToLower(namePattern) != namePattern {
		return false
	}
	// Define a regular expression pattern to match valid RFC-1035-compatible names
	pattern := regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`)
	return pattern.MatchString(namePattern)
}
