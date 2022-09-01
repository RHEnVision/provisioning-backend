package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/lzap/dejq"
)

type LaunchInstanceAWSTaskArgs struct {
	// Associated reservation
	ReservationID int64 `json:"reservation_id"`

	// Associated account
	AccountID int64 `json:"account_id"`

	// Associated public key
	PubkeyID int64 `json:"pubkey_id"`

	// Detail information
	Detail *models.AWSDetail `json:"detail"`

	// AWS AMI as fetched from image builder
	AMI string `json:"ami"`

	// The ARN fetched from Sources which is linked to a specific source
	ARN string `json:"arn"`
}

// Unmarshall arguments and handle error
func HandleLaunchInstanceAWS(ctx context.Context, job dejq.Job) error {
	args := LaunchInstanceAWSTaskArgs{}
	err := decodeJob(ctx, job, &args)
	if err != nil {
		return err
	}
	ctx = contextLogger(ctx, job.Type(), args, args.AccountID, args.ReservationID)

	jobErr := handleLaunchInstanceAWS(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
	return jobErr
}

// Job logic, when error is returned the job status is updated accordingly
func handleLaunchInstanceAWS(ctx context.Context, args *LaunchInstanceAWSTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started launch instance AWS job")

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing launch instance AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)

	client, _ := clients.GetEC2Client(ctx)
	stsClient, err := clients.GetSTSClient(ctx)
	if err != nil {
		return fmt.Errorf("cannot initialize sts client: %w", err)
	}

	crd, err := stsClient.AssumeRole(args.ARN)
	if err != nil {
		return fmt.Errorf("cannot assume role: %w", err)
	}

	newEC2Client, err := client.CreateEC2ClientFromConfig(crd)
	if err != nil {
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}

	pkD, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		return fmt.Errorf("cannot get pubkey dao: %w", err)
	}

	pk, err := pkD.GetById(ctx, args.PubkeyID)
	if err != nil {
		return fmt.Errorf("cannot get pubkey by id: %w", err)
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

	logger.Trace().Msg("Executing RunInstances")
	instances, awsReservationId, err := newEC2Client.RunInstances(ctx, args.Detail.Name, args.Detail.Amount, types.InstanceType(args.Detail.InstanceType), args.AMI, pk.Name, userData)
	if err != nil {
		return fmt.Errorf("cannot run instances: %w", err)
	}

	resD, err := dao.GetReservationDao(ctx)
	if err != nil {
		return fmt.Errorf("cannot GetReservationDao: %w", err)
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
