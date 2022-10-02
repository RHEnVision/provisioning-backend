package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func init() {
	dao.GetAccountDao = getAccountDao
}

type accountDao struct{}

func getAccountDao(ctx context.Context) (dao.AccountDao, error) {
	return &accountDao{}, nil
}

func (x *accountDao) Create(ctx context.Context, account *models.Account) error {
	query := `INSERT INTO accounts (account_number, org_id) VALUES ($1, $2) RETURNING id`
	err := db.Pool.QueryRow(ctx, query, account.AccountNumber, account.OrgID).Scan(&account.ID)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}

	return nil
}

func (x *accountDao) GetById(ctx context.Context, id int64) (*models.Account, error) {
	query := `SELECT * FROM accounts WHERE id = $1 LIMIT 1`
	result := &models.Account{}

	err := pgxscan.Get(ctx, db.Pool, result, query, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

// GetOrCreateByIdentity is not in a single transaction because it's used heavily in each request.
// This can result in duplicate error if two requests try this at once in which case an error is returned
// leading to 500 HTTP failure. Since caching of accounts was recently added, we might consider rewriting
// this function into a single transaction at some point.
func (x *accountDao) GetOrCreateByIdentity(ctx context.Context, orgId string, accountNumber string) (*models.Account, error) {
	logger := ctxval.Logger(ctx)

	// Try to find by org ID first
	acc, err := x.GetByOrgId(ctx, orgId)
	if err == nil {
		// Found it
		logger.Trace().Msgf("Account found via org id: %s", orgId)
		return acc, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// An error that is not "no rows" was returned
		return nil, err
	}

	// Previous search returned "no rows" error, try with account number
	acc, err = x.GetByAccountNumber(ctx, accountNumber)
	if err == nil {
		// Found it
		logger.Trace().Msgf("Account found via account number: %s", accountNumber)
		return acc, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// An error that is not "no rows" was returned
		return nil, err
	}

	acc = &models.Account{OrgID: orgId, AccountNumber: sql.NullString{String: accountNumber, Valid: accountNumber != ""}}
	if err := x.Create(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (x *accountDao) GetByAccountNumber(ctx context.Context, number string) (*models.Account, error) {
	query := `SELECT * FROM accounts WHERE account_number = $1 LIMIT 1`
	result := &models.Account{}

	err := pgxscan.Get(ctx, db.Pool, result, query, number)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *accountDao) GetByOrgId(ctx context.Context, orgId string) (*models.Account, error) {
	query := `SELECT * FROM accounts WHERE org_id = $1 LIMIT 1`
	result := &models.Account{}

	err := pgxscan.Get(ctx, db.Pool, result, query, orgId)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *accountDao) List(ctx context.Context, limit, offset int64) ([]*models.Account, error) {
	query := `SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2`
	var result []*models.Account

	rows, err := db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}
