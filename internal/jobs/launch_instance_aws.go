package jobs

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/rs/zerolog"
)

type LaunchInstanceAWSTaskArgs struct {
	// Associated reservation
	ReservationID int64

	// Region to provision the instances into
	Region string

	// Associated public key
	PubkeyID int64

	// Source ID that was used to get the ARN
	SourceID string

	// Detail information
	Detail *models.AWSDetail

	// AWS AMI as fetched from image builder
	AMI string

	// LaunchTemplateID or empty string when no template in use
	LaunchTemplateID string

	// The ARN fetched from Sources which is linked to a specific source
	ARN *clients.Authentication
}

// HandleLaunchInstanceAWS unmarshalls arguments and handles error
func HandleLaunchInstanceAWS(ctx context.Context, job *worker.Job) {
	logger := zerolog.Ctx(ctx)
	if job == nil {
		logger.Error().Msg("No job for HandleLaunchInstanceAWS")
		return
	}

	args, ok := job.Args.(LaunchInstanceAWSTaskArgs)
	if !ok {
		err := fmt.Errorf("%w: job %s, reservation: %#v", ErrTypeAssertion, job.ID, job.Args)
		logger.Error().Err(err).Msg("Type assertion error for job")
		return
	}

	// ensure panic finishes the job
	logger = ptr.To(logger.With().Int64("reservation_id", args.ReservationID).Logger())
	ctx = logger.WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("%w: %s", ErrPanicInJob, r)
			finishWithError(ctx, args.ReservationID, panicErr)
		}
	}()

	jobErr := DoEnsurePubkeyOnAWS(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
		return
	}

	jobErr = DoLaunchInstanceAWS(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
		return
	}

	jobErr = FetchInstancesDescriptionAWS(ctx, &args)
	if jobErr != nil {
		finishWithError(ctx, args.ReservationID, jobErr)
	}

	finishJob(ctx, args.ReservationID, jobErr)
}

// DoEnsurePubkeyOnAWS is a job logic, when error is returned the job status is updated accordingly
func DoEnsurePubkeyOnAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "DoEnsurePubkeyOnAWS")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Started pubkey upload AWS job")

	logger.Info().Interface("args", args).Msg("Processing pubkey upload AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Uploading public key")
	defer updateStatusAfter(ctx, args.ReservationID, "Uploaded public key", 1)

	pkDao := dao.GetPubkeyDao(ctx)
	resDao := dao.GetReservationDao(ctx)
	awsReservation, err := resDao.GetAWSById(ctx, args.ReservationID)
	if err != nil {
		span.SetStatus(codes.Error, "cannot get aws reservation by id")
		return fmt.Errorf("cannot get aws reservation by id: %w", err)
	}

	pubkey, err := pkDao.GetById(ctx, args.PubkeyID)
	if err != nil {
		span.SetStatus(codes.Error, "cannot get aws key")
		return fmt.Errorf("cannot get aws pubkey: %w", err)
	}

	// Fetch our DB record for the resource to update if necessary
	pkr, errDao := pkDao.UnscopedGetResourceBySourceAndRegion(ctx, args.PubkeyID, args.SourceID, args.Region)
	if errDao != nil {
		if errors.Is(errDao, dao.ErrNoRows) {
			pkr = &models.PubkeyResource{
				PubkeyID: pubkey.ID,
				Provider: models.ProviderTypeAWS,
				SourceID: args.SourceID,
				Region:   args.Region,
			}
		} else {
			span.SetStatus(codes.Error, "unable to check pubkey resource")
			return fmt.Errorf("unable to check pubkey resource: %w", errDao)
		}
	}

	ec2Client, err := clients.GetEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create new ec2 client from config")
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}

	// check presence on AWS first
	fingerprint := pubkey.FindAwsFingerprint(ctx)
	ec2Name, err := ec2Client.GetPubkeyName(ctx, fingerprint)
	if err != nil {
		span.SetStatus(codes.Error, "key error")

		// if not found on AWS, import
		if errors.Is(err, http.PubkeyNotFoundErr) {
			pkr.Tag = ""
			pkr.RandomizeTag()
			pkr.Handle, err = ec2Client.ImportPubkey(ctx, pubkey, pkr.FormattedTag())

			if errors.Is(err, http.DuplicatePubkeyErr) {
				// key not found by fingerprint but importing failed for duplicate err so fingerprints do not match
				return fmt.Errorf("key with fingerprint %s not found on AWS, but importing the key failed: %w", pubkey.Fingerprint, err)
			} else if err != nil {
				return fmt.Errorf("cannot upload aws pubkey: %w", err)
			}
			ec2Name = pubkey.Name
		} else {
			logger.Error().Err(err).Str("pubkey_fingerprint", fingerprint).Msg("Cannot fetch name of pubkey by its fingerprint")
			return fmt.Errorf("cannot fetch name of pubkey by its fingerprint: %w", err)
		}
	} else {
		logger.Debug().Msgf("Found pubkey by fingerprint (%s) with name '%s'", fingerprint, ec2Name)
	}

	// update the AWS key name in reservation details
	awsReservation.Detail.PubkeyName = ec2Name
	err = resDao.UnscopedUpdateAWSDetail(ctx, awsReservation.Reservation.ID, awsReservation.Detail)
	if err != nil {
		span.SetStatus(codes.Error, "failed to save AWS pubkey name to DB")
		return fmt.Errorf("failed to save AWS pubkey name to DB: %w", err)
	}

	if pkr.ID == 0 {
		err = pkDao.UnscopedCreateResource(ctx, pkr)
		if err != nil {
			span.SetStatus(codes.Error, "cannot create resource for aws pubkey")
			return fmt.Errorf("cannot create resource for aws pubkey: %w", err)
		}
	}

	return nilUnlessTimeout(ctx)
}

func DoLaunchInstanceAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "DoLaunchInstanceAWS")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Started launch instance AWS job")

	logger.Info().Interface("args", args).Msg("Processing launch instance AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)

	resD := dao.GetReservationDao(ctx)

	reservation, err := resD.GetAWSById(ctx, args.ReservationID)
	if err != nil {
		span.SetStatus(codes.Error, "cannot get aws reservation by id")
		return fmt.Errorf("cannot get aws reservation by id: %w", err)
	}

	// Generate user data
	userDataInput := userdata.UserData{
		Type:         models.ProviderTypeAWS,
		PowerOff:     args.Detail.PowerOff,
		InsightsTags: true,
	}
	userData, err := userdata.GenerateUserData(ctx, &userDataInput)
	if err != nil {
		span.SetStatus(codes.Error, "cannot generate user data")
		return fmt.Errorf("cannot generate user data: %w", err)
	}

	ec2Client, err := clients.GetEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create new ec2 client from config")
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}

	req := &clients.AWSInstanceParams{
		LaunchTemplateID: args.LaunchTemplateID,
		InstanceType:     types.InstanceType(args.Detail.InstanceType),
		AMI:              args.AMI,
		KeyName:          reservation.Detail.PubkeyName,
		UserData:         userData,
	}

	logger.Trace().Msg("Executing RunInstances")
	instances, awsReservationId, err := ec2Client.RunInstances(ctx, req, args.Detail.Amount, args.Detail.Name, reservation)
	if err != nil {
		span.SetStatus(codes.Error, "cannot run instances")
		return fmt.Errorf("cannot run instances: %w", err)
	}

	// For each instance that was created in AWS, add it as a DB record
	for _, instanceId := range instances {
		err = resD.CreateInstance(ctx, &models.ReservationInstance{
			ReservationID: args.ReservationID,
			InstanceID:    *instanceId,
		})
		if err != nil {
			span.SetStatus(codes.Error, "cannot create instance reservation")
			return fmt.Errorf("cannot create instance reservation for id %d: %w", instanceId, err)
		}
		logger.Info().Str("instance_id", *instanceId).Msgf("Created new instance via AWS reservation %s", *awsReservationId)
	}

	logger.Info().Str("aws_reservation_id", *awsReservationId).Msg("Adding aws reservation id")
	// Save the AWS reservation id in aws_reservation_details table
	err = resD.UpdateReservationIDForAWS(ctx, args.ReservationID, *awsReservationId)
	if err != nil {
		span.SetStatus(codes.Error, "cannot UpdateReservationIDForAWS")
		return fmt.Errorf("cannot UpdateReservationIDForAWS: %w", err)
	}

	return nilUnlessTimeout(ctx)
}

func FetchInstancesDescriptionAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "FetchInstancesDescriptionAWS")
	defer span.End()

	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("Started fetch instances description")

	updateStatusBefore(ctx, args.ReservationID, "Fetching instance(s) description")
	defer updateStatusAfter(ctx, args.ReservationID, "Instance(s) description fetched", 1)

	rDao := dao.GetReservationDao(ctx)
	instances, err := rDao.ListInstances(ctx, args.ReservationID)
	if err != nil {
		span.SetStatus(codes.Error, "cannot get instances list")
		return fmt.Errorf("cannot get instances list: %w", err)
	}
	instancesIDList := make([]string, len(instances))
	for i, instance := range instances {
		instancesIDList[i] = instance.InstanceID
	}
	ec2Client, err := clients.GetEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
		span.SetStatus(codes.Error, "cannot create new ec2 client from config")
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}

	err = waitAndRetry(ctx, func() error {
		instancesDescriptionList, errRetry := ec2Client.DescribeInstanceDetails(ctx, instancesIDList)
		if errRetry != nil {
			span.SetStatus(codes.Error, "cannot get list instances description")
			return fmt.Errorf("cannot get list instances description: %w", errRetry)
		}

		if len(instancesDescriptionList) == 0 {
			return ErrTryAgain
		}

		for _, instance := range instancesDescriptionList {
			errRetry := rDao.UpdateReservationInstance(ctx, args.ReservationID, instance)
			if errRetry != nil {
				span.SetStatus(codes.Error, "cannot update instance description")
				return fmt.Errorf("cannot update instance description: %w", errRetry)
			}
		}
		return nil
	}, 1000, 500, 500, 1000, 2000)

	if err != nil {
		span.SetStatus(codes.Error, "giving up")
		return fmt.Errorf("giving up: %w", err)
	}

	return nil
}
