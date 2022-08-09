package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

func AddPubkey(ctx context.Context, pubkey *models.Pubkey) error {
	pubkeyDao, err := getPubkeyDaoStub(ctx)
	if err != nil {
		return err
	}
	return pubkeyDao.Create(ctx, pubkey)
}
