package sqlx

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createReservation       = `INSERT INTO reservations (provider, account_id, status) VALUES ($1, $2, $3) RETURNING *`
	createAwsDetail         = `INSERT INTO aws_reservation_details (reservation_id, pubkey_id, source_id, instance_type, amount, image_id) VALUES ($1, $2, $3, $4, $5, $6)`
	updateReservationStatus = `UPDATE reservations SET status = $2 WHERE id = $1 RETURNING *`
	finishReservationStatus = `UPDATE reservations SET status = $2, success = $3, finished_at = now() WHERE id = $1 RETURNING *`
	deleteReservationById   = `DELETE FROM reservations WHERE id = $1`
	listReservations        = `SELECT * FROM reservations ORDER BY id LIMIT $1 OFFSET $2`
)

type reservationDaoSqlx struct {
	name            string
	create          *sqlx.Stmt
	createAwsDetail *sqlx.Stmt
	updateStatus    *sqlx.Stmt
	finishStatus    *sqlx.Stmt
	deleteById      *sqlx.Stmt
	list            *sqlx.Stmt
}

func getReservationDao(ctx context.Context) (dao.ReservationDao, error) {
	var err error
	daoImpl := reservationDaoSqlx{}
	daoImpl.name = "reservation"

	daoImpl.create, err = db.DB.PreparexContext(ctx, createReservation)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createReservation, err)
	}
	daoImpl.createAwsDetail, err = db.DB.PreparexContext(ctx, createAwsDetail)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createAwsDetail, err)
	}
	daoImpl.updateStatus, err = db.DB.PreparexContext(ctx, updateReservationStatus)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, updateReservationStatus, err)
	}
	daoImpl.finishStatus, err = db.DB.PreparexContext(ctx, finishReservationStatus)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, finishReservationStatus, err)
	}
	daoImpl.list, err = db.DB.PreparexContext(ctx, listReservations)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, listReservations, err)
	}
	daoImpl.deleteById, err = db.DB.PreparexContext(ctx, deleteReservationById)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, deleteReservationById, err)
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

	err := stmt.GetContext(ctx, reservation, reservation.Provider, reservation.AccountID, reservation.Status)
	if err != nil {
		return NewGetError(ctx, di, query, err)
	}
	return nil
}

func (di *reservationDaoSqlx) CreateAWS(ctx context.Context, reservation *models.AWSReservation) error {
	err := dao.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		query := createReservation
		stmt := di.create
		err := stmt.GetContext(ctx, reservation, reservation.Provider, reservation.AccountID, reservation.Status)
		if err != nil {
			return NewGetError(ctx, di, query, err)
		}

		query = createAwsDetail
		stmt = di.createAwsDetail
		res, err := stmt.ExecContext(ctx, reservation.ID, reservation.PubkeyID, reservation.SourceID, reservation.InstanceType, reservation.Amount, reservation.ImageID)
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

func (di *reservationDaoSqlx) List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error) {
	query := listReservations
	stmt := di.list
	var result []*models.Reservation

	err := stmt.SelectContext(ctx, &result, limit, offset)
	if err != nil {
		return nil, NewSelectError(ctx, di, query, err)
	}
	return result, nil
}

func (di *reservationDaoSqlx) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := updateReservationStatus
	stmt := di.updateStatus

	res, err := stmt.ExecContext(ctx, id, status)
	if err != nil {
		return NewExecUpdateError(ctx, di, query, err)
	}
	if rows, _ := res.RowsAffected(); rows != 1 {
		return NewUpdateMismatchAffectedError(ctx, di, 1, rows)

	}
	return nil
}

func (di *reservationDaoSqlx) Finish(ctx context.Context, id int64, success bool, status string) error {
	query := finishReservationStatus
	stmt := di.finishStatus

	res, err := stmt.ExecContext(ctx, id, status, success)
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
