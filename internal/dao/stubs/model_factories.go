package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

func AddPubkey(ctx context.Context, pubkey *models.Pubkey) error {
	pubkeyDao := getPubkeyDaoStub(ctx)
	return pubkeyDao.Create(ctx, pubkey)
}
