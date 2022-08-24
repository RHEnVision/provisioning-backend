package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	"github.com/RHEnVision/provisioning-backend/internal/clients/sts"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/userdata"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/lzap/dejq"
)

type LaunchInstanceAWSTaskArgs struct {
	AccountID     int64 `json:"account_id"`
	PubkeyID      int64 `json:"pubkey_id"`
	ReservationID int64 `json:"reservation_id"`
	// Optional name, can be blank
	Name string `json:"name"`
	// AWS AMI
	AMI string `json:"ami"`
	// Amount of instances to launch
	Amount int32 `json:"amount"`
	// Immediately power off the system
	PowerOff bool `json:"poweroff"`
	// Amazon EC2 Instance Type
	InstanceType string `json:"instance_type"`
	// The ARN fetched from Sources which is linked to a specific source
	ARN string `json:"arn"`
}

func HandleLaunchInstanceAWS(ctx context.Context, job dejq.Job) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started launch instance AWS job")

	args := LaunchInstanceAWSTaskArgs{}
	err := job.Decode(&args)
	if err != nil {
		ctxLogger.Error().Err(err).Msg("unable to decode arguments")
		return fmt.Errorf("unable to decode args: %w", err)
	}

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing launch instance AWS job")

	client := ec2.NewEC2Client(ctx)
	stsClient, err := sts.NewSTSClient(ctx)
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
		PowerOff: args.PowerOff,
	}
	userData, err := userdata.GenerateUserData(&userDataInput)
	if err != nil {
		return fmt.Errorf("cannot generate user data: %w", err)
	}
	logger.Trace().Bool("userdata", true).Msg(string(userData))

	logger.Info().Msg("Starting running instances on AWS")
	instances, awsReservationId, err := newEC2Client.RunInstances(ctx, args.Name, args.Amount, types.InstanceType(args.InstanceType), args.AMI, pk.Name, userData)
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

	// mark the reservation as finished
	rDao, err := dao.GetReservationDao(ctx)
	if err != nil {
		return fmt.Errorf("cannot get reservation DAO: %w", err)
	}
	err = rDao.UpdateStatus(ctx, args.ReservationID, "Finished")
	if err != nil {
		return fmt.Errorf("cannot update reservation status: %w", err)
	}

	return nil
}
