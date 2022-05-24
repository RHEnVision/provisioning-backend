package dao

import (
	"context"
	"database/sql"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

// Transaction represents both database connection pool and a single transaction/connection
type Transaction interface {
	// ExecContext from database/sql library
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// PrepareContext from database/sql library
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// QueryContext from database/sql library
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// QueryRowContext from database/sql library
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// QueryxContext from sqlx library
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	// QueryRowxContext from sqlx library
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	// PreparexContext from sqlx library
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	// GetContext from sqlx library
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// SelectContext from sqlx library
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// A TxFn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(Transaction) error

var GetAccountDao func(ctx context.Context, tx Transaction) (AccountDao, error)

type AccountDao interface {
	GetById(ctx context.Context, id uint64) (*models.Account, error)
	GetByAccountNumber(ctx context.Context, number string) (*models.Account, error)
	GetByOrgId(ctx context.Context, orgId string) (*models.Account, error)
	List(ctx context.Context, limit, offset uint64) ([]*models.Account, error)
}

var GetPubkeyDao func(ctx context.Context, tx Transaction) (PubkeyDao, error)

type PubkeyDao interface {
	Create(ctx context.Context, pk *models.Pubkey) error
	Update(ctx context.Context, pk *models.Pubkey) error
	GetById(ctx context.Context, id uint64) (*models.Pubkey, error)
	List(ctx context.Context, limit, offset uint64) ([]*models.Pubkey, error)
	Delete(ctx context.Context, id uint64) error
}

var GetPubkeyResourceDao func(ctx context.Context, tx Transaction) (PubkeyResourceDao, error)

type PubkeyResourceDao interface {
	Create(ctx context.Context, pkr *models.PubkeyResource) error
	Update(ctx context.Context, pkr *models.PubkeyResource) error
	Delete(ctx context.Context, id uint64) error
}
