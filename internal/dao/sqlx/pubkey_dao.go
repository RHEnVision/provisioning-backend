package sqlx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createPubkey     = `INSERT INTO pubkeys (account_id, name, body) VALUES ($1, $2, $3) RETURNING id`
	updatePubkey     = `UPDATE pubkeys SET account_id = $2, name = $3, body = $4 WHERE id = $1`
	getPubkeyById    = `SELECT * FROM pubkeys WHERE id = $1 LIMIT 1`
	deletePubkeyById = `DELETE FROM pubkeys WHERE id = $1`
	listPubkeys      = `SELECT * FROM pubkeys ORDER BY id LIMIT $1 OFFSET $2`
)

type pubkeyDaoSqlx struct {
	name       string
	create     *sqlx.Stmt
	update     *sqlx.Stmt
	getById    *sqlx.Stmt
	deleteById *sqlx.Stmt
	list       *sqlx.Stmt
}

func getPubkeyDao(ctx context.Context, tx dao.Transaction) (dao.PubkeyDao, error) {
	var err error
	daoImpl := pubkeyDaoSqlx{}
	daoImpl.name = "pubkey"

	daoImpl.create, err = tx.PreparexContext(ctx, createPubkey)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createPubkey, err)
	}
	daoImpl.update, err = tx.PreparexContext(ctx, updatePubkey)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updatePubkey, err)
	}
	daoImpl.getById, err = tx.PreparexContext(ctx, getPubkeyById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getAccountById, err)
	}
	daoImpl.list, err = tx.PreparexContext(ctx, listPubkeys)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listPubkeys, err)
	}
	daoImpl.deleteById, err = tx.PreparexContext(ctx, deletePubkeyById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deletePubkeyById, err)
	}

	return &daoImpl, nil
}

func (dao *pubkeyDaoSqlx) NameForError() string {
	return dao.name
}

func init() {
	dao.GetPubkeyDao = getPubkeyDao
}

func (dao *pubkeyDaoSqlx) Create(ctx context.Context, pubkey *models.Pubkey) error {
	query := createPubkey
	stmt := dao.create

	err := stmt.GetContext(ctx, pubkey, pubkey.AccountID, pubkey.Name, pubkey.Body)
	if err != nil {
		return NewGetError(ctx, dao, query, err)
	}
	return nil
}

func (dao *pubkeyDaoSqlx) GetById(ctx context.Context, id uint64) (*models.Pubkey, error) {
	query := getPubkeyById
	stmt := dao.getById
	result := &models.Pubkey{}

	err := stmt.GetContext(ctx, result, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewNoRowsError(ctx, dao, query)
		} else {
			return nil, NewGetError(ctx, dao, query, err)
		}
	}
	return result, nil
}

func (dao *pubkeyDaoSqlx) Update(ctx context.Context, pubkey *models.Pubkey) error {
	query := createPubkey
	stmt := dao.create

	res, err := stmt.ExecContext(ctx, pubkey.ID, pubkey.AccountID, pubkey.Name, pubkey.Body)
	if err != nil {
		return NewExecUpdateError(ctx, dao, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, dao, 1, rows)
	}
	return nil
}

func (dao *pubkeyDaoSqlx) List(ctx context.Context, limit, offset uint64) ([]*models.Pubkey, error) {
	query := listPubkeys
	stmt := dao.list
	var result []*models.Pubkey

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, dao, query, err)
	}
	return result, nil
}

func (dao *pubkeyDaoSqlx) Delete(ctx context.Context, id uint64) error {
	query := deletePubkeyById
	stmt := dao.deleteById

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return NewExecDeleteError(ctx, dao, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewDeleteMismatchAffectedError(ctx, dao, 1, rows)

	}
	return nil
}
