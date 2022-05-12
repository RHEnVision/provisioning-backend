package sqlx

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/jmoiron/sqlx"
)

const (
	createReservation     = `INSERT INTO reservations (account_id, pubkey_id, status) VALUES ($1, $2, $3) RETURNING *`
	deleteReservationById = `DELETE FROM reservations WHERE id = $1`
)

type reservationDaoSqlx struct {
	name       string
	create     *sqlx.Stmt
	deleteById *sqlx.Stmt
}

func getReservationDao(ctx context.Context) (dao.ReservationDao, error) {
	var err error
	daoImpl := reservationDaoSqlx{}
	daoImpl.name = "reservation"

	daoImpl.create, err = db.DB.PreparexContext(ctx, createReservation)
	if err != nil {
		return nil, NewPrepareStatementError(ctx, &daoImpl, createReservation, err)
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

func (di *reservationDaoSqlx) Create(ctx context.Context, reservation *models.Reservation) error {
	query := createReservation
	stmt := di.create

	err := stmt.GetContext(ctx, reservation, reservation.AccountID, reservation.PubkeyID, reservation.Status)
	if err != nil {
		return NewGetError(ctx, di, query, err)
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
