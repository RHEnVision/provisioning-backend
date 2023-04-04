package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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

// Unmarshall arguments and handle error
func HandleLaunchInstanceAWS(ctx context.Context, job *worker.Job) {
	args, ok := job.Args.(LaunchInstanceAWSTaskArgs)
	if !ok {
		err := fmt.Errorf("%w: job %s, reservation: %#v", ErrTypeAssertion, job.ID, job.Args)
		ctxval.Logger(ctx).Error().Err(err).Msg("Type assertion error for job")
		return
	}

	logger := ctxval.Logger(ctx).With().Int64("reservation_id", args.ReservationID).Logger()
	ctx = ctxval.WithLogger(ctx, &logger)

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

	finishJob(ctx, args.ReservationID, jobErr)
}

// Job logic, when error is returned the job status is updated accordingly
func DoEnsurePubkeyOnAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Started pubkey upload AWS job")

	logger.Info().Interface("args", args).Msg("Processing pubkey upload AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Uploading public key")
	defer updateStatusAfter(ctx, args.ReservationID, "Uploaded public key", 1)

	pkDao := dao.GetPubkeyDao(ctx)
	resDao := dao.GetReservationDao(ctx)
	awsReservation, err := resDao.GetAWSById(ctx, args.ReservationID)
	if err != nil {
		return fmt.Errorf("cannot get aws reservation by id: %w", err)
	}

	pubkey, err := pkDao.GetById(ctx, args.PubkeyID)
	if err != nil {
		return fmt.Errorf("cannot upload aws pubkey: %w", err)
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
			return fmt.Errorf("unable to check pubkey resource: %w", errDao)
		}
	}

	ec2Client, err := clients.GetEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}

	// check presence on AWS first
	fingerprint := pubkey.FindAwsFingerprint(ctx)
	ec2Name, err := ec2Client.GetPubkeyName(ctx, fingerprint)
	if err != nil {
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
			return fmt.Errorf("cannot fetch name of pubkey with fingerpring (%s): %w", fingerprint, err)
		}
	} else {
		logger.Debug().Msgf("Found pubkey by fingerprint (%s) with name '%s'", fingerprint, ec2Name)
	}

	// update the AWS key name in reservation details
	awsReservation.Detail.PubkeyName = ec2Name
	err = resDao.UnscopedUpdateAWSDetail(ctx, awsReservation.Reservation.ID, awsReservation.Detail)
	if err != nil {
		return fmt.Errorf("failed to save AWS pubkey name to DB: %w", err)
	}

	if pkr.ID == 0 {
		err = pkDao.UnscopedCreateResource(ctx, pkr)
		if err != nil {
			return fmt.Errorf("cannot create resource for aws pubkey: %w", err)
		}
	}

	return nilUnlessTimeout(ctx)
}

func DoLaunchInstanceAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Started launch instance AWS job")

	logger.Info().Interface("args", args).Msg("Processing launch instance AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)

	resD := dao.GetReservationDao(ctx)

	reservation, err := resD.GetAWSById(ctx, args.ReservationID)
	if err != nil {
		return fmt.Errorf("cannot get aws reservation by id: %w", err)
	}

	// Generate user data
	userDataInput := userdata.UserData{
		Type:         models.ProviderTypeAWS,
		PowerOff:     args.Detail.PowerOff,
		InsightsTags: true,
	}
	userData, err := userdata.GenerateUserData(&userDataInput)
	if err != nil {
		return fmt.Errorf("cannot generate user data: %w", err)
	}
	logger.Trace().Bool("userdata", true).Msg(string(userData))

	ec2Client, err := clients.GetEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
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
	instances, awsReservationId, err := ec2Client.RunInstances(ctx, req, args.Detail.Amount, args.Detail.Name)
	if err != nil {
		return fmt.Errorf("cannot run instances: %w", err)
	}

	// For each instance that was created in AWS, add it as a DB record
	for _, instanceId := range instances {
		err = resD.CreateInstance(ctx, &models.ReservationInstance{
			ReservationID: args.ReservationID,
			InstanceID:    *instanceId,
		})
		if err != nil {
			return fmt.Errorf("cannot create instance reservation for id %d: %w", instanceId, err)
		}
		logger.Info().Str("instance_id", *instanceId).Msgf("Created new instance via AWS reservation %s", *awsReservationId)
	}

	logger.Info().Str("aws_reservation_id", *awsReservationId).Msg("Adding aws reservation id")
	// Save the AWS reservation id in aws_reservation_details table
	err = resD.UpdateReservationIDForAWS(ctx, args.ReservationID, *awsReservationId)
	if err != nil {
		return fmt.Errorf("cannot UpdateReservationIDForAWS: %w", err)
	}

	return nilUnlessTimeout(ctx)
}

func FetchInstancesDescriptionAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Started fetch instances description")

	updateStatusBefore(ctx, args.ReservationID, "Fetching instance(s) description")
	defer updateStatusAfter(ctx, args.ReservationID, "Instance(s) description fetched", 1)

	rDao := dao.GetReservationDao(ctx)
	instances, err := rDao.ListInstances(ctx, args.ReservationID)
	if err != nil {
		return fmt.Errorf("cannot get instances list: %w", err)
	}
	instancesIDList := make([]string, len(instances))
	for i, instance := range instances {
		instancesIDList[i] = instance.InstanceID
	}
	ec2Client, err := clients.GetEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}
	backoffInterval := [...]int64{1000, 500, 500, 1000, 2000}
	for _, interval := range backoffInterval {
		time.Sleep(time.Duration(interval) * time.Millisecond)
		instancesDescriptionList, err := ec2Client.DescribeInstanceDetails(ctx, instancesIDList)
		if err != nil {
			return fmt.Errorf("cannot get list instances description: %w", err)
		}
		if len(instancesDescriptionList) > 0 {
			for _, instance := range instancesDescriptionList {
				err := rDao.UpdateReservationInstance(ctx, args.ReservationID, instance)
				if err != nil {
					return fmt.Errorf("cannot update instance description: %w", err)
				}
			}
			break
		}
	}

	return nil
}
