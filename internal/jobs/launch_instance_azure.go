package jobs

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
)

const (
	DefaultAzureResourceGroupName = "redhat-deployed"
	location                      = "eastus"
	DefaultVMName                 = "redhat-vm"
)

var LaunchInstanceAzureSteps = []string{"Prepare resource group", "Launch instance(s)"}

type LaunchInstanceAzureTaskArgs struct {
	// Associated reservation
	ReservationID int64

	// Location to provision the instances into when blank, uses the Resource Group location
	Location string

	// Associated public key
	PubkeyID int64

	// SourceID that was used to get the Subscription
	SourceID string

	// AzureImageID as fetched from image builder
	AzureImageID string

	// ResourceGroupName passed by a user, if left blank, defaults to 'redhat-deployed'
	ResourceGroupName string

	// The Subscription fetched from Sources which is linked to a specific source
	Subscription *clients.Authentication

	// The Name is used as prefix for a final name, for uniqueness we add uuid suffix to each instance name
	Name string
}

func HandleLaunchInstanceAzure(ctx context.Context, job *worker.Job) {
	logger := zerolog.Ctx(ctx)
	if job == nil {
		logger.Error().Msg("No job for HandleLaunchInstanceAzure")
		return
	}

	args, ok := job.Args.(LaunchInstanceAzureTaskArgs)
	if !ok {
		err := fmt.Errorf("%w: job %s, reservation: %#v", ErrTypeAssertion, job.ID, job.Args)
		logger.Error().Err(err).Msg("Type assertion error for job")
		return
	}

	// context and logger
	ctx, logger = reservationContextLogger(ctx, args.ReservationID)
	logger.Info().Msg("Started launch instance Azure job")

	if args.ResourceGroupName == "" {
		logger.Debug().Msg("Resource group has not been set, defaulting to 'redhat-deployed'")
		args.ResourceGroupName = DefaultAzureResourceGroupName
	}
	if args.Location != "" {
		// Match for availability zone suffix and removes it.
		// For backwards compatibility, we accept both forms and adjust here,
		// but we want to accept only region without the zone info in the future.
		res, e := regexp.MatchString(`_[1-6]\z`, args.Location)
		if e == nil && res {
			args.Location = args.Location[0 : len(args.Location)-2]
		}
	}
	if args.Name == "" {
		args.Name = DefaultVMName
	}

	// ensure panic finishes the job
	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("%w: %s", ErrPanicInJob, r)
			finishWithError(ctx, args.ReservationID, panicErr)
		}
	}()

	ctx, span := telemetry.StartSpan(ctx, "LaunchInstanceAzureJob")
	defer span.End()

	jobErr := DoEnsureAzureResourceGroup(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
		return
	}

	jobErr = DoLaunchInstanceAzure(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
		return
	}

	finishJob(ctx, args.ReservationID, jobErr)
}

func DoEnsureAzureResourceGroup(ctx context.Context, args *LaunchInstanceAzureTaskArgs) error {
	ctx, span := telemetry.StartSpan(ctx, "EnsureAzureResourceGroupStep")
	defer span.End()

	logger := zerolog.Ctx(ctx)

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Ensuring resource group presence")
	defer updateStatusAfter(ctx, args.ReservationID, "Ensured resource group presence", 1)

	azureClient, err := clients.GetAzureClient(ctx, args.Subscription)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create new Azure client")
		logger.Error().Err(err).Msg("Cannot create new Azure client")
		return fmt.Errorf("cannot create new Azure client: %w", err)
	}

	resourceGroup, err := azureClient.EnsureResourceGroup(ctx, args.ResourceGroupName, location)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create resource group")
		logger.Error().Err(err).Msg("Cannot create resource group")
		return fmt.Errorf("failed to ensure resource group: %w", err)
	}
	logger.Trace().Msgf("Using resource group id=%s", resourceGroup)

	if args.Location == "" {
		logger.Debug().Str("azure_location", resourceGroup.Location).Msg("Using location from Resource Group")
		args.Location = resourceGroup.Location
	}

	return nil
}

func DoLaunchInstanceAzure(ctx context.Context, args *LaunchInstanceAzureTaskArgs) error {
	ctx, span := telemetry.StartSpan(ctx, "LaunchInstanceAzureStep")
	defer span.End()

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
	userData, err := userdata.GenerateUserData(ctx, &userDataInput)
	if err != nil {
		span.SetStatus(codes.Error, "cannot generate user data")
		return fmt.Errorf("cannot generate user data: %w", err)
	}

	vmParams := clients.AzureInstanceParams{
		Location:          args.Location,
		ResourceGroupName: args.ResourceGroupName,
		ImageID:           args.AzureImageID,
		Pubkey:            pubkey,
		InstanceType:      clients.InstanceTypeName(reservation.Detail.InstanceSize),
		UserData:          userData,
		Tags: map[string]*string{
			"rh-rid": ptr.To(config.EnvironmentPrefix("r", strconv.FormatInt(reservation.ID, 10))),
			"rh-org": ptr.To(identity.Identity(ctx).Identity.OrgID),
		},
	}

	instanceDescriptions, err := azureClient.CreateVMs(ctx, vmParams, reservation.Detail.Amount, args.Name)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create instances")
		return fmt.Errorf("cannot create Azure instance: %w", err)
	}

	for _, instanceDescription := range instanceDescriptions {
		err = resDao.CreateInstance(ctx, &models.ReservationInstance{
			ReservationID: args.ReservationID,
			InstanceID:    instanceDescription.ID,
			Detail: models.ReservationInstanceDetail{
				PublicIPv4:  instanceDescription.IPv4,
				PrivateIPv4: instanceDescription.PrivateIPv4,
			},
		})
		if err != nil {
			span.SetStatus(codes.Error, "failed to save instance to DB")
			return fmt.Errorf("cannot create instance reservation for id %s: %w", instanceDescription.ID, err)
		}
	}

	return nil
}
