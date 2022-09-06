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
	createPubkeyResource            = `INSERT INTO pubkey_resources (pubkey_id, provider, source_id, handle, tag) VALUES ($1, $2, $3, $4, $5) RETURNING id, tag`
	deletePubkeyResourceById        = `DELETE FROM pubkey_resources WHERE id = $1`
	listByPubkeyId                  = `SELECT * FROM pubkey_resources WHERE pubkey_id = $1`
	getPubkeyResourceByProviderType = `SELECT * FROM pubkey_resources WHERE pubkey_id = $1 AND provider = $2`
)

type pubkeyResourceDaoSqlx struct {
	name              string
	create            *sqlx.Stmt
	deleteById        *sqlx.Stmt
	getByProviderType *sqlx.Stmt
	listByPubkeyId    *sqlx.Stmt
}

func getPubkeyResourceDao(ctx context.Context) (dao.PubkeyResourceDao, error) {
	var err error
	daoImpl := pubkeyResourceDaoSqlx{}
	daoImpl.name = "pubkeyResource"

	daoImpl.getByProviderType, err = db.DB.PreparexContext(ctx, getPubkeyResourceByProviderType)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getAccountById, err)
	}
	daoImpl.create, err = db.DB.PreparexContext(ctx, createPubkeyResource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createPubkeyResource, err)
	}
	daoImpl.deleteById, err = db.DB.PreparexContext(ctx, deletePubkeyResourceById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deletePubkeyResourceById, err)
	}
	daoImpl.listByPubkeyId, err = db.DB.PreparexContext(ctx, listByPubkeyId)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listByPubkeyId, err)
	}

	return &daoImpl, nil
}

func (di *pubkeyResourceDaoSqlx) NameForError() string {
	return di.name
}

func init() {
	dao.GetPubkeyResourceDao = getPubkeyResourceDao
}

func (di *pubkeyResourceDaoSqlx) UnscopedGetResourceByProviderType(ctx context.Context, pubkeyId int64, provider models.ProviderType) (*models.PubkeyResource, error) {
	query := getPubkeyResourceByProviderType
	stmt := di.getByProviderType
	result := &models.PubkeyResource{}

	err := stmt.GetContext(ctx, result, pubkeyId, provider)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewNoRowsError(ctx, di, query, err)
		} else {
			return nil, NewGetError(ctx, di, query, err)
		}
	}
	return result, nil
}

func (di *pubkeyResourceDaoSqlx) UnscopedCreate(ctx context.Context, pkr *models.PubkeyResource) error {
	query := createPubkeyResource
	stmt := di.create

	err := stmt.GetContext(ctx, pkr, pkr.PubkeyID, pkr.Provider, pkr.SourceID, pkr.Handle, pkr.Tag)
	if err != nil {
		return NewGetError(ctx, di, query, err)
	}
	return nil
}

func (di *pubkeyResourceDaoSqlx) UnscopedDelete(ctx context.Context, id int64) error {
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

func (di *pubkeyResourceDaoSqlx) UnscopedListByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error) {
	query := listByPubkeyId
	stmt := di.listByPubkeyId
	var result []*models.PubkeyResource

	err := stmt.SelectContext(ctx, &result, pkId)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}
