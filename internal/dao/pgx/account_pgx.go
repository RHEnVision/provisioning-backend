package pgx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type accountDao struct{}

func getAccountDao(ctx context.Context) dao.AccountDao {
	return &accountDao{}
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

// GetOrCreateByIdentity can be called multiple times on concurrent requests of new user accounts.
// Parameter accountNumber is stored as NULL when empty.
func (x *accountDao) GetOrCreateByIdentity(ctx context.Context, orgId string, accountNumber string) (*models.Account, error) {
	result := &models.Account{}
	account := sql.NullString{
		String: accountNumber,
		Valid:  accountNumber != "",
	}

	// Step 1: try to find the record by org or account (to prevent sequence gaps)
	// Step 2: create new record ignoring errors
	// Step 3: find the inserted record
	// Response is limited by 1, so it stops execution once at least one row is found.
	query := `
		WITH insert_and_select AS (
			INSERT INTO accounts (org_id, account_number) VALUES ($1, $2)
			ON CONFLICT DO NOTHING RETURNING *
		)
		SELECT * FROM accounts WHERE org_id=$1
		UNION
		SELECT * FROM accounts WHERE $2 IS NOT NULL AND account_number=$2
		UNION
		SELECT * FROM insert_and_select
		UNION
		SELECT * FROM accounts WHERE org_id=$1
		UNION
		SELECT * FROM accounts WHERE $2 IS NOT NULL AND account_number=$2
		LIMIT 1;`

	err := pgxscan.Get(ctx, db.Pool, result, query, orgId, account)
	if errors.Is(err, pgx.ErrNoRows) {
		// Step 4: requery in case transaction isolation. This will happen in case of simultaneous transactions when
		// the other transaction inserts the record after transaction is opened, therefore the conflict will trigger
		// and no account will be actually retrieved from the database. The driver returns no rows error in this case.
		zerolog.Ctx(ctx).Warn().Bool("requery", true).Msgf("Organization id %s account requery", orgId)
		err = pgxscan.Get(ctx, db.Pool, result, query, orgId, account)
		if err != nil {
			return nil, fmt.Errorf("pgx requery error: %w", err)
		}
	} else if err != nil {
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
