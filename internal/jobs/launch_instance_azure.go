package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

type LaunchInstanceAzureTaskArgs struct {
	// Associated reservation
	ReservationID int64

	// Location to provision the instances into
	Location string

	// Associated public key
	PubkeyID int64

	// SourceID that was used to get the Subscription
	SourceID string

	// AzureImageID as fetched from image builder
	AzureImageID string

	// The Subscription fetched from Sources which is linked to a specific source
	Subscription *clients.Authentication
}

func HandleLaunchInstanceAzure(ctx context.Context, job *worker.Job) {
	args, ok := job.Args.(LaunchInstanceAzureTaskArgs)
	if !ok {
		ctxval.Logger(ctx).Error().Msgf("Type assertion error for job %s, unable to finish reservation: %#v", job.ID, job.Args)
		return
	}

	logger := ctxval.Logger(ctx).With().Int64("reservation_id", args.ReservationID).Logger()
	ctx = ctxval.WithLogger(ctx, &logger)

	jobErr := DoLaunchInstanceAzure(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
}

func DoLaunchInstanceAzure(ctx context.Context, args *LaunchInstanceAzureTaskArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Started launch instance Azure job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)

	pkDao := dao.GetPubkeyDao(ctx)
	resDao := dao.GetReservationDao(ctx)

	pubkey, err := pkDao.GetById(ctx, args.PubkeyID)
	if err != nil {
		return fmt.Errorf("cannot get public key by id: %w", err)
	}

	reservation, err := resDao.GetAzureById(ctx, args.ReservationID)
	if err != nil {
		return fmt.Errorf("cannot get azure reservation by id: %w", err)
	}

	azureClient, err := clients.GetAzureClient(ctx, args.Subscription)
	if err != nil {
		return fmt.Errorf("cannot create new Azure client: %w", err)
	}

	// TODO create multiple
	instanceId, err := azureClient.CreateVM(ctx, args.AzureImageID, pubkey, clients.InstanceTypeName(reservation.Detail.InstanceSize))
	if err != nil {
		return fmt.Errorf("cannot create Azure instance(s): %w", err)
	}
	err = resDao.CreateInstance(ctx, &models.ReservationInstance{
		ReservationID: args.ReservationID,
		InstanceID:    *instanceId,
	})
	if err != nil {
		return fmt.Errorf("cannot create instance reservation for id %d: %w", instanceId, err)
	}
	logger.Debug().Msgf("Created new instance (%s) via Azure CreateVM", *instanceId)

	logger.Debug().Msg("Finished launch instance Azure job")

	return nil
}
