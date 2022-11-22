package pgx

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

func init() {
	dao.GetPubkeyDao = getPubkeyDao
}

type pubkeyDao struct{}

func getPubkeyDao(ctx context.Context) dao.PubkeyDao {
	return &pubkeyDao{}
}

func (x *pubkeyDao) validate(ctx context.Context, pubkey *models.Pubkey) error {
	if vError := models.Validate(ctx, pubkey); vError != nil {
		return fmt.Errorf("validate: %w", vError)
	}
	if tError := models.Transform(ctx, pubkey); tError != nil {
		return fmt.Errorf("transform: %w", tError)
	}

	return nil
}

func (x *pubkeyDao) Create(ctx context.Context, pubkey *models.Pubkey) error {
	query := `INSERT INTO pubkeys (account_id, name, body, fingerprint) VALUES ($1, $2, $3, $4) RETURNING id`

	pubkey.AccountID = ctxval.AccountId(ctx)

	if vError := x.validate(ctx, pubkey); vError != nil {
		return fmt.Errorf("pubkey validation: %w", vError)
	}

	err := db.Pool.QueryRow(ctx, query, pubkey.AccountID, pubkey.Name, pubkey.Body, pubkey.Fingerprint).Scan(&pubkey.ID)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}

	return nil
}

func (x *pubkeyDao) GetById(ctx context.Context, id int64) (*models.Pubkey, error) {
	query := `SELECT * FROM pubkeys WHERE account_id = $1 AND id = $2 LIMIT 1`
	accountId := ctxval.AccountId(ctx)
	result := &models.Pubkey{}

	err := pgxscan.Get(ctx, db.Pool, result, query, accountId, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *pubkeyDao) Update(ctx context.Context, pubkey *models.Pubkey) error {
	query := `UPDATE pubkeys SET name = $3, body = $4, fingerprint = $5 WHERE account_id = $1 AND id = $2`
	accountId := ctxval.AccountId(ctx)

	if vError := x.validate(ctx, pubkey); vError != nil {
		return fmt.Errorf("pubkey validation: %w", vError)
	}

	tag, err := db.Pool.Exec(ctx, query, accountId, pubkey.ID, pubkey.Name, pubkey.Body, pubkey.Fingerprint)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *pubkeyDao) List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error) {
	query := `SELECT * FROM pubkeys WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	accountId := ctxval.AccountId(ctx)
	var result []*models.Pubkey

	rows, err := db.Pool.Query(ctx, query, accountId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *pubkeyDao) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pubkeys WHERE account_id = $1 AND id = $2`
	accountId := ctxval.AccountId(ctx)

	tag, err := db.Pool.Exec(ctx, query, accountId, id)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *pubkeyDao) UnscopedCreateResource(ctx context.Context, pkr *models.PubkeyResource) error {
	query := `INSERT INTO pubkey_resources
    	(pubkey_id, provider, source_id, handle, tag, region)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, tag`

	err := db.Pool.QueryRow(ctx, query,
		pkr.PubkeyID,
		pkr.Provider,
		pkr.SourceID,
		pkr.Handle,
		pkr.Tag,
		pkr.Region).Scan(&pkr.ID, &pkr.Tag)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}

	return nil
}

func (x *pubkeyDao) UnscopedGetResourceBySourceAndRegion(ctx context.Context, pubkeyId int64, sourceId string, region string) (*models.PubkeyResource, error) {
	query := `SELECT * FROM pubkey_resources WHERE pubkey_id = $1 AND source_id = $2 AND region = $3`
	result := &models.PubkeyResource{}

	err := pgxscan.Get(ctx, db.Pool, result, query, pubkeyId, sourceId, region)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *pubkeyDao) UnscopedListResourcesByPubkeyId(ctx context.Context, id int64) ([]*models.PubkeyResource, error) {
	query := `SELECT * FROM pubkey_resources WHERE pubkey_id = $1`
	var result []*models.PubkeyResource

	rows, err := db.Pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *pubkeyDao) UnscopedDeleteResource(ctx context.Context, id int64) error {
	query := `DELETE FROM pubkey_resources WHERE id = $1`

	tag, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}
