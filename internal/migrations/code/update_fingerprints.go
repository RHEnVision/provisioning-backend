package code

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
)

// UpdateFingerprints calls appropriate DAO function, see the DAO interface for docs.
func UpdateFingerprints(ctx context.Context) error {
	pkd := dao.GetServiceDao(ctx)
	count, err := pkd.RecalculatePubkeyFingerprints(ctx)
	if err != nil {
		return fmt.Errorf("error when updating fingerprints: %w", err)
	}
	ctxval.Logger(ctx).Info().Msgf("Total number of updated pubkey records: %d", count)
	return nil
}
