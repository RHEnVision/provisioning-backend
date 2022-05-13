package sqlx

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	getById            = `SELECT * FROM accounts WHERE id = $1 LIMIT 1`
	getByAccountNumber = `SELECT * FROM accounts WHERE account_number = $1 LIMIT 1`
	getByOrgId         = `SELECT * FROM accounts WHERE org_id = $1 LIMIT 1`
	list               = `SELECT * FROM accounts ORDER BY id LIMIT $1 OFFSET $2`
)

type accountDaoSqlx struct {
	getById            *sqlx.Stmt
	getByAccountNumber *sqlx.Stmt
	getByOrgId         *sqlx.Stmt
	list               *sqlx.Stmt
}

func getAccountDao(ctx context.Context) (dao.AccountDao, error) {
	var err error
	daoImpl := accountDaoSqlx{}

	daoImpl.getById, err = db.DB.PreparexContext(ctx, getById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, getById, err)
	}
	daoImpl.getByAccountNumber, err = db.DB.PreparexContext(ctx, getByAccountNumber)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, list, err)
	}
	daoImpl.getByOrgId, err = db.DB.PreparexContext(ctx, getByOrgId)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, list, err)
	}
	daoImpl.list, err = db.DB.PreparexContext(ctx, list)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, list, err)
	}

	return &daoImpl, nil
}

func init() {
	dao.GetAccountDao = getAccountDao
}

func (a *accountDaoSqlx) GetById(ctx context.Context, id uint64) (*models.Account, error) {
	result := &models.Account{}
	err := a.getById.GetContext(ctx, result, id)
	if err != nil {
		return nil, NewGetError(ctx, "get by id", err)
	}
	return result, nil
}

func (a *accountDaoSqlx) GetByAccountNumber(ctx context.Context, number string) (*models.Account, error) {
	result := &models.Account{}
	err := a.getByAccountNumber.GetContext(ctx, result, number)
	if err != nil {
		return nil, NewGetError(ctx, "get by id", err)
	}
	return result, nil
}

func (a *accountDaoSqlx) GetByOrgId(ctx context.Context, orgId string) (*models.Account, error) {
	result := &models.Account{}
	err := a.getByOrgId.GetContext(ctx, result, orgId)
	if err != nil {
		return nil, NewGetError(ctx, "get by id", err)
	}
	return result, nil
}

func (a *accountDaoSqlx) List(ctx context.Context, limit, offset uint64) ([]*models.Account, error) {
	var result []*models.Account
	err := a.list.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewGetError(ctx, "list", err)
	}
	return result, nil
}
