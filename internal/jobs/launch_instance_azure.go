package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const (
	resourceGroupName = "redhat-deployed"
	location          = "eastus"
	vmNamePrefix      = "redhat-vm"
)

var LaunchInstanceAzureSteps = []string{"Prepare resource group", "Launch instance(s)"}

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
		err := fmt.Errorf("%w: job %s, reservation: %#v", ErrTypeAssertion, job.ID, job.Args)
		ctxval.Logger(ctx).Error().Err(err).Msg("Type assertion error for job")
		return
	}

	logger := ctxval.Logger(ctx).With().Int64("reservation_id", args.ReservationID).Logger()
	ctx = ctxval.WithLogger(ctx, &logger)

	logger.Info().Msg("Started launch instance Azure job")
	ctx, span := otel.Tracer(TraceName).Start(ctx, "LaunchInstanceAzureJob")
	defer span.End()

	jobErr := DoEnsureAzureResourceGroup(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
		return
	}

	jobErr = DoLaunchInstanceAzure(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)

	logger.Info().Msg("Finished launch instance Azure job")
}

func DoEnsureAzureResourceGroup(ctx context.Context, args *LaunchInstanceAzureTaskArgs) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "EnsureAzureResourceGroupStep")
	defer span.End()

	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Ensuring resource group presence")
	defer updateStatusAfter(ctx, args.ReservationID, "Ensured resource group presence", 1)

	azureClient, err := clients.GetAzureClient(ctx, args.Subscription)
	if err != nil {
		return fmt.Errorf("cannot create new Azure client: %w", err)
	}

	resourceGroupID, err := azureClient.EnsureResourceGroup(ctx, resourceGroupName, location)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create resource group")
		logger.Error().Err(err).Msg("Cannot create resource group")
		return fmt.Errorf("failed to ensure resource group: %w", err)
	}
	logger.Trace().Msgf("Using resource group id=%s", *resourceGroupID)
	return nil
}

func DoLaunchInstanceAzure(ctx context.Context, args *LaunchInstanceAzureTaskArgs) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "LaunchInstanceAzureStep")
	defer span.End()

	logger := ctxval.Logger(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)

	pkDao := dao.GetPubkeyDao(ctx)
	resDao := dao.GetReservationDao(ctx)

	pubkey, err := pkDao.GetById(ctx, args.PubkeyID)
	if err != nil {
		span.SetStatus(codes.Error, "cannot get public key by id")
		return fmt.Errorf("cannot get public key by id: %w", err)
	}

	reservation, err := resDao.GetAzureById(ctx, args.ReservationID)
	if err != nil {
		span.SetStatus(codes.Error, "cannot get azure reservation record")
		return fmt.Errorf("cannot get azure reservation by id: %w", err)
	}

	azureClient, err := clients.GetAzureClient(ctx, args.Subscription)
	if err != nil {
		span.SetStatus(codes.Error, "cannot instantiate Azure client")
		return fmt.Errorf("failed to instantiate Azure client: %w", err)
	}
	// Generate user data
	userDataInput := userdata.UserData{
		Type:         models.ProviderTypeAzure,
		PowerOff:     reservation.Detail.PowerOff,
		InsightsTags: true,
	}
	userData, err := userdata.GenerateUserData(&userDataInput)
	if err != nil {
		return fmt.Errorf("cannot generate user data: %w", err)
	}
	logger.Trace().Bool("userdata", true).Msg(string(userData))

	vmParams := clients.AzureInstanceParams{
		Location:          location,
		ResourceGroupName: resourceGroupName,
		ImageID:           args.AzureImageID,
		Pubkey:            pubkey,
		InstanceType:      clients.InstanceTypeName(reservation.Detail.InstanceSize),
		UserData:          userData,
	}

	instanceDescriptions, err := azureClient.CreateVMs(ctx, vmParams, reservation.Detail.Amount, vmNamePrefix)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create instances")
		return fmt.Errorf("cannot create Azure instance: %w", err)
	}

	for _, instanceDescription := range instanceDescriptions {
		err = resDao.CreateInstance(ctx, &models.ReservationInstance{
			ReservationID: args.ReservationID,
			InstanceID:    instanceDescription.ID,
			Detail: models.ReservationInstanceDetail{
				PublicIPv4: instanceDescription.PublicIPv4,
			},
		})
		if err != nil {
			span.SetStatus(codes.Error, "failed to save instance to DB")
			return fmt.Errorf("cannot create instance reservation for id %s: %w", instanceDescription.ID, err)
		}
	}

	return nil
}
