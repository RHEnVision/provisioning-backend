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
	createSource          = `INSERT INTO sources (account_id, source_id, name, auth_id) VALUES ($1, $2, $3, $4) RETURNING id,created_at`
	getSourceById         = `SELECT * FROM sources WHERE id = $1;`
	listSources           = `SELECT * FROM sources ORDER BY id LIMIT $1 OFFSET $2`
	deleteSourceById      = `DELETE FROM sources WHERE id = $1 RETURNING source_id`
	getSourcesByAccountId = `SELECT * FROM sources WHERE account_id = $1 ORDER BY id LIMIT $1 OFFSET $2;`
)

type sourceDaoSqlx struct {
	name           string
	getById        *sqlx.Stmt
	list           *sqlx.Stmt
	create         *sqlx.Stmt
	getByAccountId *sqlx.Stmt
	deleteById     *sqlx.Stmt
}

func getSourceDao(ctx context.Context) (dao.SourceDao, error) {
	var err error
	daoImpl := sourceDaoSqlx{}
	daoImpl.name = "source"

	daoImpl.getById, err = db.DB.PreparexContext(ctx, getSourceById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getSourceById, err)
	}
	daoImpl.list, err = db.DB.PreparexContext(ctx, listSources)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listSources, err)
	}
	daoImpl.create, err = db.DB.PreparexContext(ctx, createSource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createSource, err)
	}
	daoImpl.getByAccountId, err = db.DB.PreparexContext(ctx, getSourcesByAccountId)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getSourcesByAccountId, err)
	}
	return &daoImpl, nil
}

func (di *sourceDaoSqlx) NameForError() string {
	return di.name
}

func init() {
	dao.GetSourceDao = getSourceDao
}

func (di *sourceDaoSqlx) Create(ctx context.Context, source *models.Source) (*models.Source, error) {
	query := createSource
	stmt := di.create
	err := stmt.GetContext(ctx, source, source.AccountID, source.SourceID, source.Name, source.AuthID)
	if err != nil {
		return nil, NewGetError(ctx, di, query, err)
	}
	return source, nil
}

func (di *sourceDaoSqlx) GetById(ctx context.Context, id uint64) (*models.Source, error) {
	query := getSourceById
	stmt := di.getById
	result := &models.Source{}

	err := stmt.GetContext(ctx, result, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NewNoRowsError(ctx, di, query)
		} else {
			return nil, NewGetError(ctx, di, query, err)
		}
	}
	return result, nil
}

func (di *sourceDaoSqlx) List(ctx context.Context, limit uint64, offset uint64) ([]*models.Source, error) {
	query := listSources
	stmt := di.list
	var result []*models.Source

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *sourceDaoSqlx) GetByAccountId(ctx context.Context, accountId uint64, limit uint64, offset uint64) ([]*models.Source, error) {
	query := getSourcesByAccountId
	stmt := di.getByAccountId
	var result []*models.Source

	err := stmt.SelectContext(ctx, &result, accountId, limit, offset)
	if err != nil {
		return nil, NewGetError(ctx, di, query, err)
	}
	return result, nil
}

func (di *sourceDaoSqlx) Delete(ctx context.Context, id uint64) (*models.Source, error) {
	query := deleteSourceById
	stmt := di.deleteById
	result := &models.Source{}

	err := stmt.GetContext(ctx, result, id)
	if err != nil {
		return nil, NewGetError(ctx, di, query, err)
	}
	return result, nil

}
