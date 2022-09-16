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
	createPubkey                    = `INSERT INTO pubkeys (account_id, name, body, fingerprint) VALUES ($1, $2, $3, $4) RETURNING id`
	updatePubkey                    = `UPDATE pubkeys SET name = $3, body = $4 WHERE account_id = $1 AND id = $2`
	getPubkeyById                   = `SELECT * FROM pubkeys WHERE account_id = $1 AND id = $2 LIMIT 1`
	deletePubkeyById                = `DELETE FROM pubkeys WHERE account_id = $1 AND id = $2`
	listPubkeys                     = `SELECT * FROM pubkeys WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	createPubkeyResource            = `INSERT INTO pubkey_resources (pubkey_id, provider, source_id, handle, tag, region) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, tag`
	getPubkeyResourceByProviderType = `SELECT * FROM pubkey_resources WHERE pubkey_id = $1 AND provider = $2`
	listPubkeyResourceByPubkeyId    = `SELECT * FROM pubkey_resources WHERE pubkey_id = $1`
	deletePubkeyResourceById        = `DELETE FROM pubkey_resources WHERE id = $1`
)

type pubkeyDaoSqlx struct {
	name                            string
	create                          *sqlx.Stmt
	update                          *sqlx.Stmt
	getById                         *sqlx.Stmt
	deleteById                      *sqlx.Stmt
	list                            *sqlx.Stmt
	createPubkeyResource            *sqlx.Stmt
	deletePubkeyResourceById        *sqlx.Stmt
	getPubkeyResourceByProviderType *sqlx.Stmt
	listPubkeyResourceByPubkeyId    *sqlx.Stmt
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
	daoImpl.createPubkeyResource, err = db.DB.PreparexContext(ctx, createPubkeyResource)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createPubkeyResource, err)
	}
	daoImpl.getPubkeyResourceByProviderType, err = db.DB.PreparexContext(ctx, getPubkeyResourceByProviderType)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getPubkeyResourceByProviderType, err)
	}
	daoImpl.listPubkeyResourceByPubkeyId, err = db.DB.PreparexContext(ctx, listPubkeyResourceByPubkeyId)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listPubkeyResourceByPubkeyId, err)
	}
	daoImpl.deletePubkeyResourceById, err = db.DB.PreparexContext(ctx, deletePubkeyResourceById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deletePubkeyResourceById, err)
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
	if validationErr := models.Validate(ctx, pubkey); validationErr != nil {
		return dao.NewValidationError(ctx, di, pubkey, validationErr)
	}
	if err := models.Transform(ctx, pubkey); err != nil {
		return dao.TransformationError{
			Message: "cannot generate SSH fingerprint",
			Context: ctx,
			Err:     err,
		}
	}

	query := createPubkey
	stmt := di.create

	err := stmt.GetContext(ctx, pubkey, ctxAccountId(ctx), pubkey.Name, pubkey.Body, pubkey.Fingerprint)
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

func (di *pubkeyDaoSqlx) UnscopedCreate(ctx context.Context, pkr *models.PubkeyResource) error {
	query := createPubkeyResource
	stmt := di.createPubkeyResource

	err := stmt.GetContext(ctx, pkr, pkr.PubkeyID, pkr.Provider, pkr.SourceID, pkr.Handle, pkr.Tag, pkr.Region)
	if err != nil {
		return NewGetError(ctx, di, query, err)
	}
	return nil
}

func (di *pubkeyDaoSqlx) UnscopedGetResourceByProviderType(ctx context.Context, pubkeyId int64, provider models.ProviderType) (*models.PubkeyResource, error) {
	query := getPubkeyResourceByProviderType
	stmt := di.getPubkeyResourceByProviderType
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

func (di *pubkeyDaoSqlx) UnscopedListByPubkeyId(ctx context.Context, pkId int64) ([]*models.PubkeyResource, error) {
	query := listPubkeyResourceByPubkeyId
	stmt := di.listPubkeyResourceByPubkeyId
	var result []*models.PubkeyResource

	err := stmt.SelectContext(ctx, &result, pkId)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *pubkeyDaoSqlx) UnscopedDelete(ctx context.Context, id int64) error {
	query := deletePubkeyResourceById
	stmt := di.deletePubkeyResourceById

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return NewExecDeleteError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewDeleteMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}
