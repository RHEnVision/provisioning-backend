package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	"github.com/RHEnVision/provisioning-backend/internal/clients/sts"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/lzap/dejq"
)

type PubkeyUploadAWSTaskArgs struct {
	AccountID     int64  `json:"account_id"`
	ReservationID int64  `json:"reservation_id"`
	PubkeyID      int64  `json:"pubkey_id"`
	SourceID      string `json:"source_id"`
	ARN           string `json:"arn"`
}

func HandlePubkeyUploadAWS(ctx context.Context, job dejq.Job) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started pubkey upload AWS job")

	args := PubkeyUploadAWSTaskArgs{}
	err := job.Decode(&args)
	if err != nil {
		ctxLogger.Error().Err(err).Msg("unable to decode arguments")
		return fmt.Errorf("unable to decode args: %w", err)
	}

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing pubkey upload AWS job")

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

	pkr.Handle, err = newEC2Client.ImportPubkey(pubkey, pkr.FormattedTag())
	if err != nil {
		if errors.Is(err, ec2.DuplicatePubkeyErr) {
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
