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

type PubkeyUploadAWSTaskArgs struct {
	AccountID     int64  `json:"account_id"`
	ReservationID int64  `json:"reservation_id"`
	Region        string `json:"region"`
	PubkeyID      int64  `json:"pubkey_id"`
	SourceID      string `json:"source_id"`
	ARN           string `json:"arn"`
}

// Unmarshall arguments and handle error
func HandlePubkeyUploadAWS(ctx context.Context, job dejq.Job) error {
	args := PubkeyUploadAWSTaskArgs{}
	err := decodeJob(ctx, job, &args)
	if err != nil {
		return err
	}
	ctx = contextLogger(ctx, job.Type(), args, args.AccountID, args.ReservationID)

	jobErr := handlePubkeyUploadAWS(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
	return jobErr
}

// Job logic, when error is returned the job status is updated accordingly
func handlePubkeyUploadAWS(ctx context.Context, args *PubkeyUploadAWSTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started pubkey upload AWS job")

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing pubkey upload AWS job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Uploading public key")
	defer updateStatusAfter(ctx, args.ReservationID, "Uploaded public key", 1)

	pkDao, err := dao.GetPubkeyDao(ctx)
	if err != nil {
		return fmt.Errorf("cannot upload aws pubkey: %w", err)
	}

	pubkey, err := pkDao.GetById(ctx, args.PubkeyID)
	if err != nil {
		return fmt.Errorf("cannot upload aws pubkey: %w", err)
	}

	pkrDao, err := dao.GetPubkeyResourceDao(ctx)
	if err != nil {
		return fmt.Errorf("cannot upload aws pubkey: %w", err)
	}

	// check presence first
	skip := true
	pkrCheck, errDao := pkrDao.GetResourceByProviderType(ctx, args.PubkeyID, models.ProviderTypeAWS)
	if errDao != nil {
		var e dao.NoRowsError
		if errors.As(errDao, &e) {
			skip = false
		} else {
			return fmt.Errorf("unable to check pubkey resource: %w", errDao)
		}
	}

	if skip {
		logger.Info().Msgf("SSH key-pair '%s' already present, no upload needed", pkrCheck.Handle)
		return nil
	}

	// create new resource with randomized tag
	pkr := models.PubkeyResource{
		PubkeyID: pubkey.ID,
		Provider: models.ProviderTypeAWS,
		SourceID: args.SourceID,
	}
	pkr.RandomizeTag()

	// upload to cloud with a tag
	ec2Client, err := clients.GetCustomerEC2Client(ctx, args.ARN, args.Region)
	if err != nil {
		return fmt.Errorf("cannot create new ec2 client from config: %w", err)
	}

	pkr.Handle, err = ec2Client.ImportPubkey(pubkey, pkr.FormattedTag())
	if err != nil {
		if errors.Is(err, http.DuplicatePubkeyErr) {
			logger.Warn().Msgf("Pubkey '%s' already present, skipping", pubkey.Name)
		} else {
			return fmt.Errorf("cannot upload aws pubkey: %w", err)
		}
	}

	// create resource with handle
	err = pkrDao.Create(ctx, &pkr)
	if err != nil {
		return fmt.Errorf("cannot upload aws pubkey: %w", err)
	}

	return nil
}
