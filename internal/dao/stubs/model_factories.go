package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

func AddPubkey(ctx context.Context, pubkey *models.Pubkey) error {
	pubkeyDao := getPubkeyDaoStub(ctx)
	return pubkeyDao.Create(ctx, pubkey)
}

func AddAWSReservation(ctx context.Context, reservation *models.AWSReservation) error {
	reservationDao := getReservationDaoStub(ctx)
	return reservationDao.CreateAWS(ctx, reservation)
}
