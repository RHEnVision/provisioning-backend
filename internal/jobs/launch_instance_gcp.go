package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/notifications"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/rs/zerolog"
)

var LaunchInstanceGCPSteps = []string{"Launch instance(s)", "Fetch instance(s) description"}

type LaunchInstanceGCPTaskArgs struct {
	// Associated reservation
	ReservationID int64

	// Zone to provision the instances into
	Zone string

	// Associated public key
	PubkeyID int64

	// Detail information
	Detail *models.GCPDetail

	// GCP image name as fetched from image builder
	ImageName string

	// The project id from Sources which is linked to a specific source
	ProjectID *clients.Authentication

	// Launch template id or empty string when no template in use
	LaunchTemplateID string
}

// HandleLaunchInstanceGCP unmarshalls arguments and handles error
func HandleLaunchInstanceGCP(ctx context.Context, job *worker.Job) {
	logger := zerolog.Ctx(ctx)
	if job == nil {
		logger.Error().Msg("No job for HandleLaunchInstanceGCP")
		return
	}
	args, ok := job.Args.(LaunchInstanceGCPTaskArgs)
	if !ok {
		err := fmt.Errorf("%w: job %s, reservation: %#v", ErrTypeAssertion, job.ID, job.Args)
		logger.Error().Err(err).Msg("Type assertion error for job")
		return
	}

	logger = ptr.To(logger.With().Int64("reservation_id", args.ReservationID).Logger())
	ctx = logger.WithContext(ctx)
	nc := notifications.GetNotificationClient(ctx)

	jobErr := DoLaunchInstanceGCP(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
		nc.FailedLaunch(ctx, args.ReservationID, jobErr)
		return
	}

	jobErr = FetchInstancesDescriptionGCP(ctx, &args)
	if jobErr != nil {
		nc.FailedLaunch(ctx, args.ReservationID, jobErr)
	} else {
		nc.SuccessfulLaunch(ctx, args.ReservationID)
	}
	finishJob(ctx, args.ReservationID, jobErr)
}

// DoLaunchInstanceGCP is a job logic, when error is returned the job status is updated accordingly
func DoLaunchInstanceGCP(ctx context.Context, args *LaunchInstanceGCPTaskArgs) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Started launch instance GCP job")
	logger.Info().Interface("args", args).Msg("Processing launch instance GCP job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)
	name := "inst-####"
	pkD := dao.GetPubkeyDao(ctx)

	pk, err := pkD.GetById(ctx, args.PubkeyID)
	if err != nil {
		return fmt.Errorf("cannot get pubkey by id: %w", err)
	}

	gcpClient, err := clients.GetGCPClient(ctx, args.ProjectID)
	if err != nil {
		return fmt.Errorf("cannot get gcp client: %w", err)
	}

	// Generate user data
	userDataInput := userdata.UserData{
		Type:         models.ProviderTypeGCP,
		PowerOff:     args.Detail.PowerOff,
		InsightsTags: true,
	}
	userData, err := userdata.GenerateUserData(ctx, &userDataInput)
	if err != nil {
		return fmt.Errorf("cannot generate user data: %w", err)
	}

	if args.Detail.NamePattern != nil {
		name = fmt.Sprintf("%s-#####", *args.Detail.NamePattern)
	}
	params := &clients.GCPInstanceParams{
		NamePattern:      ptr.To(name),
		ImageName:        args.ImageName,
		MachineType:      args.Detail.MachineType,
		Zone:             args.Zone,
		KeyBody:          pk.Body,
		StartupScript:    string(userData),
		ReservationID:    args.ReservationID,
		UUID:             args.Detail.UUID,
		LaunchTemplateID: args.LaunchTemplateID,
	}

	instances, opName, err := gcpClient.InsertInstances(ctx, params, args.Detail.Amount)
	if err != nil {
		return fmt.Errorf("cannot run instances for gcp client: %w", err)
	}

	rDao := dao.GetReservationDao(ctx)

	err = rDao.UpdateOperationNameForGCP(ctx, args.ReservationID, *opName)
	if err != nil {
		return fmt.Errorf("cannot update operation name for GCP : %w", err)
	}

	// For each instance that was created in GCP, add it as a DB record
	for _, instanceId := range instances {
		err = rDao.CreateInstance(ctx, &models.ReservationInstance{
			ReservationID: args.ReservationID,
			InstanceID:    *instanceId,
		})
		if err != nil {
			return fmt.Errorf("cannot create instance reservation for id %d: %w", instanceId, err)
		}
		logger.Info().Str("instance_id", *instanceId).Msgf("Created new instance via GCP reservation %s", *opName)
	}

	return nil
}

func FetchInstancesDescriptionGCP(ctx context.Context, args *LaunchInstanceGCPTaskArgs) error {
	logger := *zerolog.Ctx(ctx)
	logger.Debug().Msg("Started Fetch Instances Description GCP")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Fetching instances description")
	defer updateStatusAfter(ctx, args.ReservationID, "Fetched instances description", 1)

	rDao := dao.GetReservationDao(ctx)

	gcpClient, err := clients.GetGCPClient(ctx, args.ProjectID)
	if err != nil {
		return fmt.Errorf("cannot get gcp client: %w", err)
	}
	ids, err := gcpClient.ListInstancesIDsByLabel(ctx, args.Detail.UUID)
	if err != nil {
		return fmt.Errorf("cannot list instances ids by tag: %w", err)
	}

	for _, id := range ids {
		var instanceDesc *clients.InstanceDescription
		err = waitAndRetry(ctx, func() error {
			instanceDesc, err = gcpClient.GetInstanceDescriptionByID(ctx, *id, args.Zone)

			if err != nil {
				return fmt.Errorf("cannot get instance description : %w", err)
			}

			if instanceDesc.PublicIPv4 == "" {
				return ErrTryAgain
			}

			return nil
		}, 1, 500, 500, 1000, 2000, 2000)

		if err != nil {
			logger.Error().Err(err).Str("instance_id", *id).Msg("Cannot get instance description, skipping")

			// try to get the others
			continue
		}

		err = rDao.UpdateReservationInstance(ctx, args.ReservationID, instanceDesc)
		if err != nil {
			return fmt.Errorf("cannot update instance description: %w", err)
		}
	}

	return nil
}
