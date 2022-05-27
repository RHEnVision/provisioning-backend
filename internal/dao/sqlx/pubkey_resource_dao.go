package sqlx

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createPubkeyResource     = `INSERT INTO pubkey_resources (pubkey_id, provider, handle, tag) VALUES ($1, $2, $3, $4) RETURNING id, tag`
	updatePubkeyResource     = `UPDATE pubkey_resources SET pubkey_id = $2, provider = $3, handle = $4 WHERE id = $1`
	deletePubkeyResourceById = `DELETE FROM pubkey_resources WHERE id = $1`
)

type pubkeyResourceDaoSqlx struct {
	name       string
	create     *sqlx.Stmt
	update     *sqlx.Stmt
	deleteById *sqlx.Stmt
}

func getPubkeyResourceDao(ctx context.Context) (dao.PubkeyResourceDao, error) {
	var err error
	daoImpl := pubkeyResourceDaoSqlx{}
	daoImpl.name = "pubkeyResource"

	daoImpl.create, err = db.DB.PreparexContext(ctx, createPubkeyResource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createPubkeyResource, err)
	}
	daoImpl.update, err = db.DB.PreparexContext(ctx, updatePubkeyResource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updatePubkeyResource, err)
	}
	daoImpl.deleteById, err = db.DB.PreparexContext(ctx, deletePubkeyResourceById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deletePubkeyResourceById, err)
	}

	return &daoImpl, nil
}

func (di *pubkeyResourceDaoSqlx) NameForError() string {
	return di.name
}

func init() {
	dao.GetPubkeyResourceDao = getPubkeyResourceDao
}

func (di *pubkeyResourceDaoSqlx) Create(ctx context.Context, pkr *models.PubkeyResource) error {
	query := createPubkeyResource
	stmt := di.create

	err := stmt.GetContext(ctx, pkr, pkr.PubkeyID, pkr.Provider, pkr.Handle, pkr.Tag)
	if err != nil {
		return NewGetError(ctx, di, query, err)
	}
	return nil
}

func (di *pubkeyResourceDaoSqlx) Update(ctx context.Context, pkr *models.PubkeyResource) error {
	query := updatePubkeyResource
	stmt := di.update

	res, err := stmt.ExecContext(ctx, pkr.ID, pkr.PubkeyID, pkr.Provider, pkr.Handle)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *pubkeyResourceDaoSqlx) Delete(ctx context.Context, id uint64) error {
	query := deletePubkeyResourceById
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
