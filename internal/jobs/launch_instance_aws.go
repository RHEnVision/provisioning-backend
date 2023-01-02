package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type LaunchInstanceAWSTaskArgs struct {
	// Associated reservation
	ReservationID int64 `json:"reservation_id"`

	// Associated account
	AccountID int64 `json:"account_id"`

	// Region to provision the instances into
	Region string `json:"region"`

	// Associated public key
	PubkeyID int64 `json:"pubkey_id"`

	// Source ID that was used to get the ARN
	SourceID string

	// Detail information
	Detail *models.AWSDetail `json:"detail"`

	// AWS AMI as fetched from image builder
	AMI string `json:"ami"`

	// The ARN fetched from Sources which is linked to a specific source
	ARN *clients.Authentication `json:"arn"`
}

var LaunchInstanceAWSTask = queue.RegisterTask("launch_instance_aws", func(ctx context.Context, args LaunchInstanceAWSTaskArgs) error {
	ctx = prepareContext(ctx, "launch_instance_aws", args, args.AccountID, args.ReservationID)
	err := HandleEnsurePubkeyOnAWS(ctx, args)
	finishStep(ctx, args.ReservationID, err)
	if err != nil {
		return err
	}

	err = HandleLaunchInstanceAWS(ctx, args)
	finishStep(ctx, args.ReservationID, err)
	return err
})

func EnqueueLaunchInstanceAWS(ctx context.Context, args LaunchInstanceAWSTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)

	t := LaunchInstanceAWSTask.WithArgs(ctx, args)
	ctxLogger.Debug().Str("tid", t.ID).Msg("Adding launch AWS job task")
	err := queue.JobQueue.Add(t)
	if err != nil {
		return fmt.Errorf("unable to enqueue task: %w", err)
	}

	return nil
}

func HandleEnsurePubkeyOnAWS(ctx context.Context, args LaunchInstanceAWSTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started pubkey upload AWS job")

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing pubkey upload AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Uploading public key")
	defer updateStatusAfter(ctx, args.ReservationID, "Uploaded public key", 1)

	pkDao := dao.GetPubkeyDao(ctx)
	resDao := dao.GetReservationDao(ctx)
	awsReservation, err := resDao.GetAWSById(ctx, args.ReservationID)
	if err != nil {
		return fmt.Errorf("cannot upload aws pubkey: %w", err)
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
	fingerprint, err := pubkey.FingerprintAWS()
	if err != nil {
		return fmt.Errorf("failed to calculate MD5 fingerprint: %w", err)
	}
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

	// TODO this do not need to be in the DB anymore since we switched from separate jobs into a single job.
	// Update the AWS key name in reservation details:
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

	return nil
}

func HandleLaunchInstanceAWS(ctx context.Context, args LaunchInstanceAWSTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started launch instance AWS job")

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
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
		PowerOff: args.Detail.PowerOff,
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

	logger.Trace().Msg("Executing RunInstances")
	instances, awsReservationId, err := ec2Client.RunInstances(ctx, args.Detail.Name, args.Detail.Amount, types.InstanceType(args.Detail.InstanceType), args.AMI, reservation.Detail.PubkeyName, userData)
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

	return nil
}
