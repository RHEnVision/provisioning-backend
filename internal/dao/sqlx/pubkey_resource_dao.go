package sqlx

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createPubkeyResource     = `INSERT INTO pubkey_resources (pubkey_id, provider, handle) VALUES ($1, $2, $3) RETURNING id, tag`
	updatePubkeyResource     = `UPDATE pubkey_resources SET pubkey_id = $2, provider = $3, handle = $4 WHERE id = $1`
	deletePubkeyResourceById = `DELETE FROM pubkey_resources WHERE id = $1`
)

type pubkeyResourceDaoSqlx struct {
	name       string
	create     *sqlx.Stmt
	update     *sqlx.Stmt
	deleteById *sqlx.Stmt
}

func getPubkeyResourceDao(ctx context.Context, tx dao.Transaction) (dao.PubkeyResourceDao, error) {
	var err error
	daoImpl := pubkeyResourceDaoSqlx{}
	daoImpl.name = "pubkeyResource"

	daoImpl.create, err = tx.PreparexContext(ctx, createPubkeyResource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createPubkeyResource, err)
	}
	daoImpl.update, err = tx.PreparexContext(ctx, updatePubkeyResource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updatePubkeyResource, err)
	}
	daoImpl.deleteById, err = tx.PreparexContext(ctx, deletePubkeyResourceById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deletePubkeyResourceById, err)
	}

	return &daoImpl, nil
}

func (dao *pubkeyResourceDaoSqlx) NameForError() string {
	return dao.name
}

func init() {
	dao.GetPubkeyResourceDao = getPubkeyResourceDao
}

func (dao *pubkeyResourceDaoSqlx) Create(ctx context.Context, pkr *models.PubkeyResource) error {
	query := createPubkeyResource
	stmt := dao.create

	err := stmt.GetContext(ctx, pkr, pkr.PubkeyID, pkr.Provider, pkr.Handle)
	if err != nil {
		return NewGetError(ctx, dao, query, err)
	}
	return nil
}

func (dao *pubkeyResourceDaoSqlx) Update(ctx context.Context, pkr *models.PubkeyResource) error {
	query := updatePubkeyResource
	stmt := dao.update

	res, err := stmt.ExecContext(ctx, pkr.ID, pkr.PubkeyID, pkr.Provider, pkr.Handle)
	if err != nil {
		return NewExecUpdateError(ctx, dao, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, dao, 1, rows)
	}
	return nil
}

func (dao *pubkeyResourceDaoSqlx) Delete(ctx context.Context, id uint64) error {
	query := deletePubkeyResourceById
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
