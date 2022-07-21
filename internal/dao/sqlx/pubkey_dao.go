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
	query := createPubkey
	stmt := di.create

	err := stmt.GetContext(ctx, pubkey, pubkey.AccountID, pubkey.Name, pubkey.Body)
	if err != nil {
		return NewGetError(ctx, di, query, err)
	}
	return nil
}

func (di *pubkeyDaoSqlx) GetById(ctx context.Context, id int64) (*models.Pubkey, error) {
	query := getPubkeyById
	stmt := di.getById
	result := &models.Pubkey{}

	err := stmt.GetContext(ctx, result, id)
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
	query := updatePubkey
	stmt := di.update

	res, err := stmt.ExecContext(ctx, pubkey.ID, pubkey.AccountID, pubkey.Name, pubkey.Body)
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

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *pubkeyDaoSqlx) Delete(ctx context.Context, id int64) error {
	query := deletePubkeyById
	stmt := di.deleteById

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return NewExecDeleteError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewDeleteMismatchAffectedError(ctx, di, 1, rows)

	}
	return nil
}
