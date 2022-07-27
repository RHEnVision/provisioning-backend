package sqlx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createPubkey     = `INSERT INTO pubkeys (account_id, name, body) VALUES ($1, $2, $3) RETURNING id`
	updatePubkey     = `UPDATE pubkeys SET name = $3, body = $4 WHERE account_id = $1 AND id = $2`
	getPubkeyById    = `SELECT * FROM pubkeys WHERE account_id = $1 AND id = $2 LIMIT 1`
	deletePubkeyById = `DELETE FROM pubkeys WHERE account_id = $1 AND id = $2`
	listPubkeys      = `SELECT * FROM pubkeys WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
)

type pubkeyDaoSqlx struct {
	name       string
	create     *sqlx.Stmt
	update     *sqlx.Stmt
	getById    *sqlx.Stmt
	deleteById *sqlx.Stmt
	list       *sqlx.Stmt
}

func getPubkeyDao(ctx context.Context) (dao.PubkeyDao, error) {
	var err error
	daoImpl := pubkeyDaoSqlx{}
	daoImpl.name = "pubkey"

	daoImpl.create, err = db.DB.PreparexContext(ctx, createPubkey)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createPubkey, err)
	}
	daoImpl.update, err = db.DB.PreparexContext(ctx, updatePubkey)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updatePubkey, err)
	}
	daoImpl.getById, err = db.DB.PreparexContext(ctx, getPubkeyById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getPubkeyById, err)
	}
	daoImpl.list, err = db.DB.PreparexContext(ctx, listPubkeys)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listPubkeys, err)
	}
	daoImpl.deleteById, err = db.DB.PreparexContext(ctx, deletePubkeyById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deletePubkeyById, err)
	}

	return &daoImpl, nil
}

func (di *pubkeyDaoSqlx) NameForError() string {
	return di.name
}

func init() {
	dao.GetPubkeyDao = getPubkeyDao
}

func (di *pubkeyDaoSqlx) Create(ctx context.Context, pubkey *models.Pubkey) error {
	if pubkey.AccountID == 0 {
		pubkey.AccountID = ctxAccountId(ctx)
	}
	if pubkey.AccountID != ctxAccountId(ctx) {
		return dao.WrongTenantError
	}

	query := createPubkey
	stmt := di.create

	err := stmt.GetContext(ctx, pubkey, ctxAccountId(ctx), pubkey.Name, pubkey.Body)
	if err != nil {
		return NewCreateError(ctx, di, query, err)
	}
	return nil
}

func (di *pubkeyDaoSqlx) GetById(ctx context.Context, id int64) (*models.Pubkey, error) {
	query := getPubkeyById
	stmt := di.getById
	result := &models.Pubkey{}

	err := stmt.GetContext(ctx, result, ctxAccountId(ctx), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewNoRowsError(ctx, di, query, err)
		} else {
			return nil, NewGetError(ctx, di, query, err)
		}
	}
	return result, nil
}

func (di *pubkeyDaoSqlx) Update(ctx context.Context, pubkey *models.Pubkey) error {
	if pubkey.AccountID == 0 {
		pubkey.AccountID = ctxAccountId(ctx)
	}
	if pubkey.AccountID != ctxAccountId(ctx) {
		return dao.WrongTenantError
	}

	query := updatePubkey
	stmt := di.update

	res, err := stmt.ExecContext(ctx, ctxAccountId(ctx), pubkey.ID, pubkey.Name, pubkey.Body)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *pubkeyDaoSqlx) List(ctx context.Context, limit, offset int64) ([]*models.Pubkey, error) {
	query := listPubkeys
	stmt := di.list
	var result []*models.Pubkey

	err := stmt.SelectContext(ctx, &result, ctxAccountId(ctx), limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *pubkeyDaoSqlx) Delete(ctx context.Context, id int64) error {
	query := deletePubkeyById
	stmt := di.deleteById

	res, err := stmt.ExecContext(ctx, ctxAccountId(ctx), id)
	if err != nil {
		return NewExecDeleteError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewDeleteMismatchAffectedError(ctx, di, 1, rows)

	}
	return nil
}
