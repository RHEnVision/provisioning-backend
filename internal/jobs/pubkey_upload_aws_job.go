package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients/ec2"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/lzap/dejq"
)

type PubkeyUploadAWSTaskArgs struct {
	AccountID     int64 `json:"account_id"`
	ReservationID int64 `json:"reservation_id"`
	PubkeyID      int64 `json:"pubkey_id"`
}

func EnqueuePubkeyUploadAWS(ctx context.Context, args *PubkeyUploadAWSTaskArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Interface("arg", args).Msgf("Enqueuing pubkey upload AWS job: %+v", args)

	pj := dejq.PendingJob{
		Type: TypePubkeyUploadAws,
		Body: args,
	}
	err := Queue.Enqueue(ctx, pj)
	if err != nil {
		return fmt.Errorf("unable to enqueue: %w", err)
	}

	return nil
}

func HandlePubkeyUploadAWS(ctx context.Context, job dejq.Job) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Started pubkey upload AWS job")

	args := PubkeyUploadAWSTaskArgs{}
	err := job.Decode(&args)
	if err != nil {
		logger.Error().Err(err).Msg("unable to decode arguments")
		return fmt.Errorf("unable to decode args: %w", err)
	}
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

	// create new resource with randomized tag
	pkr := models.PubkeyResource{
		PubkeyID: pubkey.ID,
		Provider: models.ProviderTypeAWS,
	}
	pkr.RandomizeTag()

	// upload to cloud with a tag
	client := ec2.NewEC2Client(ctx)
	pkr.Handle, err = client.ImportPubkey(pubkey, pkr.FormattedTag())
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

	return nil
}
