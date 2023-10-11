package pgx

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/rs/zerolog"
)

func init() {
	dao.GetServiceDao = getServiceDao
}

type serviceDao struct{}

func getServiceDao(_ context.Context) dao.ServiceDao {
	return &serviceDao{}
}

func UnscopedUpdatePubkey(ctx context.Context, pubkey *models.Pubkey) error {
	query := `
		UPDATE pubkeys SET
			type = $2,
			name = $3,
			body = $4,
			fingerprint = $5,
			fingerprint_legacy = $6
		WHERE id = $1`

	tag, err := db.Pool.Exec(ctx, query, pubkey.ID, pubkey.Type, pubkey.Name, pubkey.Body, pubkey.Fingerprint, pubkey.FingerprintLegacy)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row, got %d: %w", tag.RowsAffected(), dao.ErrAffectedMismatch)
	}
	return nil
}

// RecalculatePubkeyFingerprints recalculates fingerprints for all keys which have a blank value in any of
// the fingerprints or type. The type column with value "test" is also considered as pubkey which needs
// to be recalculated as this is used in tests. Fingerprints starting with "SHA256" are also considered the same.
func (x *serviceDao) RecalculatePubkeyFingerprints(ctx context.Context) (int, error) {
	total := 0
	query := `SELECT * FROM pubkeys WHERE type = '' OR type = 'test' OR fingerprint LIKE 'SHA256:%' OR fingerprint = '' OR fingerprint_legacy = ''`
	logger := zerolog.Ctx(ctx)

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return total, fmt.Errorf("pgx error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pk models.Pubkey

		err = pgxscan.ScanRow(&pk, rows)
		if err != nil {
			return total, fmt.Errorf("pgx scan error: %w", err)
		}

		logger.Trace().Msgf("Pubkey before: %+v", pk)
		if tError := models.Transform(ctx, &pk); tError != nil {
			return total, fmt.Errorf("transform: %w", tError)
		}
		logger.Trace().Msgf("Pubkey after: %+v", pk)

		logger.Debug().Msgf("Updating pubkey fingerprints of %d named %s", pk.ID, pk.Name)
		err = UnscopedUpdatePubkey(ctx, &pk)
		if err != nil {
			return total, fmt.Errorf("pgx update error: %w", err)
		}
		total += 1
	}
	if err := rows.Err(); err != nil {
		return total, fmt.Errorf("pgx error: %w", err)
	}

	return total, nil
}
