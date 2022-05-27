package sqlx

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	getAccountById            = `SELECT * FROM accounts WHERE id = $1 LIMIT 1`
	getAccountByAccountNumber = `SELECT * FROM accounts WHERE account_number = $1 LIMIT 1`
	getAccountByOrgId         = `SELECT * FROM accounts WHERE org_id = $1 LIMIT 1`
	listAccounts              = `SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2`
)

type accountDaoSqlx struct {
	name               string
	getById            *sqlx.Stmt
	getByAccountNumber *sqlx.Stmt
	getByOrgId         *sqlx.Stmt
	list               *sqlx.Stmt
}

func getAccountDao(ctx context.Context, tx dao.Transaction) (dao.AccountDao, error) {
	var err error
	daoImpl := accountDaoSqlx{}
	daoImpl.name = "account"

	daoImpl.getById, err = tx.PreparexContext(ctx, getAccountById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getAccountById, err)
	}
	daoImpl.getByAccountNumber, err = tx.PreparexContext(ctx, getAccountByAccountNumber)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listAccounts, err)
	}
	daoImpl.getByOrgId, err = tx.PreparexContext(ctx, getAccountByOrgId)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listAccounts, err)
	}
	daoImpl.list, err = tx.PreparexContext(ctx, listAccounts)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listAccounts, err)
	}

	return &daoImpl, nil
}

func (dao *accountDaoSqlx) NameForError() string {
	return dao.name
}

func init() {
	dao.GetAccountDao = getAccountDao
}

func (dao *accountDaoSqlx) GetById(ctx context.Context, id uint64) (*models.Account, error) {
	query := getAccountById
	stmt := dao.getById
	result := &models.Account{}

	err := stmt.GetContext(ctx, result, id)
	if err != nil {
		return nil, NewGetError(ctx, dao, query, err)
	}
	return result, nil
}

func (dao *accountDaoSqlx) GetByAccountNumber(ctx context.Context, number string) (*models.Account, error) {
	query := getAccountByAccountNumber
	stmt := dao.getByAccountNumber
	result := &models.Account{}

	err := stmt.GetContext(ctx, result, number)
	if err != nil {
		return nil, NewGetError(ctx, dao, query, err)
	}
	return result, nil
}

func (dao *accountDaoSqlx) GetByOrgId(ctx context.Context, orgId string) (*models.Account, error) {
	query := getAccountByOrgId
	stmt := dao.getByOrgId
	result := &models.Account{}

	err := stmt.GetContext(ctx, result, orgId)
	if err != nil {
		return nil, NewGetError(ctx, dao, query, err)
	}
	return result, nil
}

func (dao *accountDaoSqlx) List(ctx context.Context, limit, offset uint64) ([]*models.Account, error) {
	query := listAccounts
	stmt := dao.list
	var result []*models.Account

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, dao, query, err)
	}
	return result, nil
}
