package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/lzap/dejq"
)

type EnsurePubkeyOnAWSTaskArgs struct {
	AccountID     int64                   `json:"account_id"`
	ReservationID int64                   `json:"reservation_id"`
	Region        string                  `json:"region"`
	PubkeyID      int64                   `json:"pubkey_id"`
	SourceID      string                  `json:"source_id"`
	ARN           *clients.Authentication `json:"arn"`
}

// HandleEnsurePubkeyOnAWS takes pubkey and ensures the pubkey is present on AWS in requested region.
// It saves the name of the pubkey in models.AWSReservation in models.AWSDetail.
// This only unmarshall arguments and handles error, processing function is not exported.
func HandleEnsurePubkeyOnAWS(ctx context.Context, job dejq.Job) error {
	args := EnsurePubkeyOnAWSTaskArgs{}
	err := decodeJob(ctx, job, &args)
	if err != nil {
		return err
	}
	ctx = contextLogger(ctx, job.Type(), args, args.AccountID, args.ReservationID)

	jobErr := handleEnsurePubkeyOnAWS(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
	return jobErr
}

// Job logic, when error is returned the job status is updated accordingly
func handleEnsurePubkeyOnAWS(ctx context.Context, args *EnsurePubkeyOnAWSTaskArgs) error {
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
	ec2Name, err := ec2Client.GetPubkeyName(ctx, pubkey.Fingerprint)
	if err != nil {
		// if not found on AWS, import
		if errors.Is(err, http.PubkeyNotFoundErr) {
			pkr.Tag = ""
			pkr.RandomizeTag()
			pkr.Handle, err = ec2Client.ImportPubkey(ctx, pubkey, pkr.FormattedTag())
			if err != nil {
				return fmt.Errorf("cannot upload aws pubkey: %w", err)
			}
			ec2Name = pubkey.Name
		} else {
			return fmt.Errorf("cannot fetch name of pubkey with fingerpring (%s): %w", pubkey.Fingerprint, err)
		}
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

	return nil
}
