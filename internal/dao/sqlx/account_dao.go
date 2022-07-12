package sqlx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createAccount             = `INSERT INTO accounts (account_number, org_id) VALUES ($1, $2) ON CONFLICT DO UPDATE RETURNING id`
	getAccountById            = `SELECT * FROM accounts WHERE id = $1 LIMIT 1`
	getAccountByAccountNumber = `SELECT * FROM accounts WHERE account_number = $1 LIMIT 1`
	getAccountByOrgId         = `SELECT * FROM accounts WHERE org_id = $1 LIMIT 1`
	listAccounts              = `SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2`
)

type accountDaoSqlx struct {
	name               string
	create             *sqlx.Stmt
	getById            *sqlx.Stmt
	getByAccountNumber *sqlx.Stmt
	getByOrgId         *sqlx.Stmt
	list               *sqlx.Stmt
}

func getAccountDao(ctx context.Context) (dao.AccountDao, error) {
	var err error
	daoImpl := accountDaoSqlx{}
	daoImpl.name = "account"

	daoImpl.create, err = db.DB.PreparexContext(ctx, createAccount)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createAccount, err)
	}
	daoImpl.getById, err = db.DB.PreparexContext(ctx, getAccountById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getAccountById, err)
	}
	daoImpl.getByAccountNumber, err = db.DB.PreparexContext(ctx, getAccountByAccountNumber)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listAccounts, err)
	}
	daoImpl.getByOrgId, err = db.DB.PreparexContext(ctx, getAccountByOrgId)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listAccounts, err)
	}
	daoImpl.list, err = db.DB.PreparexContext(ctx, listAccounts)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listAccounts, err)
	}

	return &daoImpl, nil
}

func (di *accountDaoSqlx) NameForError() string {
	return di.name
}

func init() {
	dao.GetAccountDao = getAccountDao
}

func (di *accountDaoSqlx) Create(ctx context.Context, account *models.Account) error {
	query := createAccount
	stmt := di.create

	err := stmt.GetContext(ctx, account, account.AccountNumber, account.OrgID)
	if err != nil {
		return NewCreateError(ctx, di, query, err)
	}
	return nil
}

func (di *accountDaoSqlx) GetById(ctx context.Context, id int64) (*models.Account, error) {
	query := getAccountById
	stmt := di.getById
	result := &models.Account{}

	err := stmt.GetContext(ctx, result, id)
	if err != nil {
		return nil, NewGetError(ctx, di, query, err)
	}
	return result, nil
}

func (di *accountDaoSqlx) GetOrCreateByIdentity(ctx context.Context, orgId string, accountNumber string) (*models.Account, error) {
	var acc, err = di.GetByOrgId(ctx, orgId)
	var norows *dao.NoRowsError
	if err == nil {
		return acc, nil
	} else if errors.As(err, &norows) {
		ctxval.Logger(ctx).Debug().Msgf("Account not found by org id %s", orgId)
	} else {
		return nil, err
	}
	acc, err = di.GetByAccountNumber(ctx, accountNumber)
	if err == nil {
		return acc, nil
	} else if errors.As(err, &norows) {
		ctxval.Logger(ctx).Debug().Msgf("Account not found by account number %s", accountNumber)
	} else {
		return nil, err
	}
	// TODO create HMSPROV-135
	return nil, err
}

func (di *accountDaoSqlx) GetByAccountNumber(ctx context.Context, number string) (*models.Account, error) {
	query := getAccountByAccountNumber
	stmt := di.getByAccountNumber
	result := &models.Account{}

	err := stmt.GetContext(ctx, result, number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewNoRowsError(ctx, di, query, err)
		} else {
			return nil, NewGetError(ctx, di, query, err)
		}
	}
	return result, nil
}

func (di *accountDaoSqlx) GetByOrgId(ctx context.Context, orgId string) (*models.Account, error) {
	query := getAccountByOrgId
	stmt := di.getByOrgId
	result := &models.Account{}

	err := stmt.GetContext(ctx, result, orgId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewNoRowsError(ctx, di, query, err)
		} else {
			return nil, NewGetError(ctx, di, query, err)
		}
	}
	return result, nil
}

func (di *accountDaoSqlx) List(ctx context.Context, limit, offset int64) ([]*models.Account, error) {
	query := listAccounts
	stmt := di.list
	var result []*models.Account

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}
