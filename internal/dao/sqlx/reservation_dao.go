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
	createReservation         = `INSERT INTO reservations (provider, account_id, steps, status) VALUES ($1, $2, $3, $4) RETURNING *`
	createAwsDetail           = `INSERT INTO aws_reservation_details (reservation_id, pubkey_id, source_id, image_id, detail) VALUES ($1, $2, $3, $4, $5)`
	createGcpDetail           = `INSERT INTO gcp_reservation_details (reservation_id, pubkey_id, source_id, image_id, detail) VALUES ($1, $2, $3, $4, $5)`
	updateReservationStatus   = `UPDATE reservations SET status = $2, step = step + $3 WHERE id = $1 RETURNING *`
	updateReservationIDForAWS = `UPDATE aws_reservation_details SET aws_reservation_id = $2 WHERE reservation_id = $1 RETURNING *`
	updateOperationNameForGCP = `UPDATE gcp_reservation_details SET gcp_operation_name = $2 WHERE reservation_id = $1 RETURNING *`
	finishReservationSuccess  = `UPDATE reservations SET success = true, finished_at = now() WHERE id = $1 RETURNING *`
	finishReservationError    = `UPDATE reservations SET success = false, error = $2, finished_at = now() WHERE id = $1 RETURNING *`
	deleteReservationById     = `DELETE FROM reservations WHERE id = $1`
	getReservationById        = `SELECT * FROM reservations WHERE account_id = $1 AND id = $2 LIMIT 1`
	listReservations          = `SELECT * FROM reservations WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3`
	createInstance            = `INSERT INTO reservation_instances (reservation_id, instance_id) VALUES ($1, $2)`
	listInstanceReservations  = `SELECT * FROM reservation_instances ORDER BY reservation_id LIMIT $1 OFFSET $2`

	getAWSReservationById = `SELECT id, provider, account_id, created_at, steps, step, status, error, finished_at, success,
    	pubkey_id, source_id, image_id, aws_reservation_id, detail
		FROM reservations, aws_reservation_details
		WHERE account_id = $1 AND id = $2 AND id = reservation_id AND provider = provider_type_aws() LIMIT 1`
)

type reservationDaoSqlx struct {
	name                      string
	create                    *sqlx.Stmt
	getById                   *sqlx.Stmt
	getAWSById                *sqlx.Stmt
	createAwsDetail           *sqlx.Stmt
	createGcpDetail           *sqlx.Stmt
	updateStatus              *sqlx.Stmt
	finishSuccess             *sqlx.Stmt
	finishError               *sqlx.Stmt
	deleteById                *sqlx.Stmt
	list                      *sqlx.Stmt
	createInstance            *sqlx.Stmt
	updateReservationIDForAWS *sqlx.Stmt
	updateOperationNameForGCP *sqlx.Stmt
	listInstanceReservations  *sqlx.Stmt
}

func getReservationDao(ctx context.Context) (dao.ReservationDao, error) {
	var err error
	daoImpl := reservationDaoSqlx{}
	daoImpl.name = "reservation"

	daoImpl.create, err = db.DB.PreparexContext(ctx, createReservation)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createReservation, err)
	}
	daoImpl.getById, err = db.DB.PreparexContext(ctx, getReservationById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getReservationById, err)
	}
	daoImpl.getAWSById, err = db.DB.PreparexContext(ctx, getAWSReservationById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, getAWSReservationById, err)
	}
	daoImpl.createAwsDetail, err = db.DB.PreparexContext(ctx, createAwsDetail)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createAwsDetail, err)
	}
	daoImpl.createGcpDetail, err = db.DB.PreparexContext(ctx, createGcpDetail)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createGcpDetail, err)
	}
	daoImpl.updateStatus, err = db.DB.PreparexContext(ctx, updateReservationStatus)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updateReservationStatus, err)
	}
	daoImpl.finishSuccess, err = db.DB.PreparexContext(ctx, finishReservationSuccess)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, finishReservationSuccess, err)
	}
	daoImpl.finishError, err = db.DB.PreparexContext(ctx, finishReservationError)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, finishReservationError, err)
	}
	daoImpl.list, err = db.DB.PreparexContext(ctx, listReservations)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listReservations, err)
	}
	daoImpl.deleteById, err = db.DB.PreparexContext(ctx, deleteReservationById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deleteReservationById, err)
	}
	daoImpl.createInstance, err = db.DB.PreparexContext(ctx, createInstance)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createInstance, err)
	}
	daoImpl.updateReservationIDForAWS, err = db.DB.PreparexContext(ctx, updateReservationIDForAWS)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updateReservationIDForAWS, err)
	}
	daoImpl.updateOperationNameForGCP, err = db.DB.PreparexContext(ctx, updateOperationNameForGCP)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updateOperationNameForGCP, err)
	}
	daoImpl.listInstanceReservations, err = db.DB.PreparexContext(ctx, listInstanceReservations)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listInstanceReservations, err)
	}
	return &daoImpl, nil
}

func (di *reservationDaoSqlx) NameForError() string {
	return di.name
}

func init() {
	dao.GetReservationDao = getReservationDao
}

func (di *reservationDaoSqlx) CreateNoop(ctx context.Context, reservation *models.NoopReservation) error {
	query := createReservation
	stmt := di.create

	err := stmt.GetContext(ctx, reservation,
		reservation.Provider,
		reservation.AccountID,
		reservation.Steps,
		reservation.Status)
	if err != nil {
		return NewGetError(ctx, di, query, err)
	}
	return nil
}

func (di *reservationDaoSqlx) CreateAWS(ctx context.Context, reservation *models.AWSReservation) error {
	err := dao.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		query := createReservation
		stmt := di.create
		err := stmt.GetContext(ctx, reservation,
			reservation.Provider,
			reservation.AccountID,
			reservation.Steps,
			reservation.Status)
		if err != nil {
			return NewGetError(ctx, di, query, err)
		}

		query = createAwsDetail
		stmt = di.createAwsDetail
		res, err := stmt.ExecContext(ctx,
			reservation.ID,
			reservation.PubkeyID,
			reservation.SourceID,
			reservation.ImageID,
			reservation.Detail)
		if err != nil {
			return NewExecUpdateError(ctx, di, query, err)
		}
		if rows, _ := res.RowsAffected(); rows != 1 {
			return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
		}
		return nil
	})
	if err != nil {
		return NewTransactionError(ctx, err)
	}
	return nil
}

func (di *reservationDaoSqlx) CreateGCP(ctx context.Context, reservation *models.GCPReservation) error {
	err := dao.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		query := createReservation
		stmt := di.create
		err := stmt.GetContext(ctx, reservation,
			reservation.Provider,
			reservation.AccountID,
			reservation.Steps,
			reservation.Status)
		if err != nil {
			return NewGetError(ctx, di, query, err)
		}

		query = createGcpDetail
		stmt = di.createGcpDetail
		res, err := stmt.ExecContext(ctx,
			reservation.ID,
			reservation.PubkeyID,
			reservation.SourceID,
			reservation.ImageID,
			reservation.Detail)
		if err != nil {
			return NewExecUpdateError(ctx, di, query, err)
		}
		if rows, _ := res.RowsAffected(); rows != 1 {
			return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
		}
		return nil
	})
	if err != nil {
		return NewTransactionError(ctx, err)
	}
	return nil
}

func (di *reservationDaoSqlx) CreateInstance(ctx context.Context, reservation *models.ReservationInstance) error {
	err := dao.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		query := createInstance
		stmt := di.createInstance
		res, err := stmt.ExecContext(ctx, reservation.ReservationID, reservation.InstanceID)
		if err != nil {
			return NewExecUpdateError(ctx, di, query, err)
		}
		if rows, _ := res.RowsAffected(); rows != 1 {
			return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
		}
		return nil
	})
	if err != nil {
		return NewTransactionError(ctx, err)
	}
	return nil
}

func (di *reservationDaoSqlx) GetById(ctx context.Context, id int64) (*models.Reservation, error) {
	query := getReservationById
	stmt := di.getById
	result := &models.Reservation{}

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

func (di *reservationDaoSqlx) GetAWSById(ctx context.Context, id int64) (*models.AWSReservation, error) {
	query := getAWSReservationById
	stmt := di.getAWSById
	result := &models.AWSReservation{}

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

func (di *reservationDaoSqlx) List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error) {
	query := listReservations
	stmt := di.list
	var result []*models.Reservation

	err := stmt.SelectContext(ctx, &result, ctxAccountId(ctx), limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *reservationDaoSqlx) ListInstances(ctx context.Context, limit, offset int64) ([]*models.ReservationInstance, error) {
	query := listInstanceReservations
	stmt := di.listInstanceReservations
	var result []*models.ReservationInstance

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *reservationDaoSqlx) UpdateStatus(ctx context.Context, id int64, status string, addSteps int32) error {
	query := updateReservationStatus
	stmt := di.updateStatus

	res, err := stmt.ExecContext(ctx, id, status, addSteps)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *reservationDaoSqlx) UpdateReservationIDForAWS(ctx context.Context, id int64, awsReservationId string) error {
	query := updateReservationIDForAWS
	stmt := di.updateReservationIDForAWS

	res, err := stmt.ExecContext(ctx, id, awsReservationId)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *reservationDaoSqlx) UpdateOperationNameForGCP(ctx context.Context, id int64, gcpOperationName string) error {
	query := updateOperationNameForGCP
	stmt := di.updateOperationNameForGCP

	res, err := stmt.ExecContext(ctx, id, gcpOperationName)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *reservationDaoSqlx) FinishWithSuccess(ctx context.Context, id int64) error {
	query := finishReservationSuccess
	stmt := di.finishSuccess

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *reservationDaoSqlx) FinishWithError(ctx context.Context, id int64, errorString string) error {
	query := finishReservationError
	stmt := di.finishError

	res, err := stmt.ExecContext(ctx, id, errorString)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)
	}
	return nil
}

func (di *reservationDaoSqlx) Delete(ctx context.Context, id int64) error {
	query := deleteReservationById
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
