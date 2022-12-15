package pgx

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func init() {
	dao.GetReservationDao = getReservationDao
}

type reservationDao struct{}

func getReservationDao(ctx context.Context) dao.ReservationDao {
	return &reservationDao{}
}

func (x *reservationDao) CreateNoop(ctx context.Context, reservation *models.NoopReservation) error {
	query := `INSERT INTO reservations (provider, account_id, steps, status, step_titles) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	reservation.AccountID = ctxval.AccountId(ctx)

	err := db.Pool.QueryRow(ctx, query,
		reservation.Provider,
		reservation.AccountID,
		reservation.Steps,
		reservation.Status,
		reservation.StepTitles).Scan(&reservation.ID, &reservation.CreatedAt)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}

	return nil
}

func (x *reservationDao) CreateAWS(ctx context.Context, reservation *models.AWSReservation) error {
	txErr := dao.WithTransaction(ctx, func(tx pgx.Tx) error {
		reservation.AccountID = ctxval.AccountId(ctx)

		reservationQuery := `INSERT INTO reservations (provider, account_id, steps, status, step_titles)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
		err := db.Pool.QueryRow(ctx, reservationQuery,
			reservation.Provider,
			reservation.AccountID,
			reservation.Steps,
			reservation.Status,
			reservation.StepTitles).Scan(&reservation.ID, &reservation.CreatedAt)
		if err != nil {
			return fmt.Errorf("pgx error: %w", err)
		}

		awsQuery := `INSERT INTO aws_reservation_details (reservation_id, pubkey_id, source_id, image_id, detail)
		VALUES ($1, $2, $3, $4, $5)`
		tag, err := db.Pool.Exec(ctx, awsQuery,
			reservation.ID,
			reservation.PubkeyID,
			reservation.SourceID,
			reservation.ImageID,
			reservation.Detail)
		if err != nil {
			return fmt.Errorf("pgx error: %w", err)
		}
		if tag.RowsAffected() != 1 {
			return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("pgx transaction error: %w", txErr)
	}
	return nil
}

func (x *reservationDao) CreateGCP(ctx context.Context, reservation *models.GCPReservation) error {
	txErr := dao.WithTransaction(ctx, func(tx pgx.Tx) error {
		reservation.AccountID = ctxval.AccountId(ctx)

		reservationQuery := `INSERT INTO reservations (provider, account_id, steps, status)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`
		err := db.Pool.QueryRow(ctx, reservationQuery,
			reservation.Provider,
			reservation.AccountID,
			reservation.Steps,
			reservation.Status).Scan(&reservation.ID, &reservation.CreatedAt)
		if err != nil {
			return fmt.Errorf("pgx error: %w", err)
		}

		gcpQuery := `INSERT INTO gcp_reservation_details (reservation_id, pubkey_id, source_id, image_id, detail)
			VALUES ($1, $2, $3, $4, $5)`
		tag, err := db.Pool.Exec(ctx, gcpQuery,
			reservation.ID,
			reservation.PubkeyID,
			reservation.SourceID,
			reservation.ImageID,
			reservation.Detail)
		if err != nil {
			return fmt.Errorf("pgx error: %w", err)
		}
		if tag.RowsAffected() != 1 {
			return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
		}

		return nil
	})

	if txErr != nil {
		return fmt.Errorf("pgx transaction error: %w", txErr)
	}
	return nil
}

func (x *reservationDao) CreateInstance(ctx context.Context, instance *models.ReservationInstance) error {
	query := `INSERT INTO reservation_instances (reservation_id, instance_id) VALUES ($1, $2)`

	tag, err := db.Pool.Exec(ctx, query,
		instance.ReservationID,
		instance.InstanceID)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}

	return nil
}

func (x *reservationDao) GetById(ctx context.Context, id int64) (*models.Reservation, error) {
	query := `SELECT * FROM reservations WHERE account_id = $1 AND id = $2 LIMIT 1`
	accountId := ctxval.AccountId(ctx)
	result := &models.Reservation{}

	err := pgxscan.Get(ctx, db.Pool, result, query, accountId, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) GetAWSById(ctx context.Context, id int64) (*models.AWSReservation, error) {
	query := `SELECT id, provider, account_id, created_at, steps, step, status, error, finished_at, success,
    	pubkey_id, source_id, image_id, aws_reservation_id, detail
		FROM reservations, aws_reservation_details
		WHERE account_id = $1 AND id = $2 AND id = reservation_id AND provider = provider_type_aws() LIMIT 1`
	accountId := ctxval.AccountId(ctx)
	result := &models.AWSReservation{}

	err := pgxscan.Get(ctx, db.Pool, result, query, accountId, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error) {
	query := `SELECT * FROM reservations WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3`

	accountId := ctxval.AccountId(ctx)
	var result []*models.Reservation

	rows, err := db.Pool.Query(ctx, query, accountId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) ListInstances(ctx context.Context, reservationId int64) ([]*models.ReservationInstance, error) {
	query := `SELECT reservation_id, instance_id FROM reservation_instances, reservations
         WHERE reservation_id = reservations.id AND account_id = $1 AND reservation_id = $2`

	accountId := ctxval.AccountId(ctx)
	var result []*models.ReservationInstance

	rows, err := db.Pool.Query(ctx, query, accountId, reservationId)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}

	err = pgxscan.ScanAll(&result, rows)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) UpdateStatus(ctx context.Context, id int64, status string, addSteps int32) error {
	query := `UPDATE reservations SET status = $2, step = step + $3 WHERE id = $1`

	tag, err := db.Pool.Exec(ctx, query, id, status, addSteps)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *reservationDao) UnscopedUpdateAWSDetail(ctx context.Context, id int64, awsDetail *models.AWSDetail) error {
	query := `UPDATE aws_reservation_details SET detail = $2 WHERE reservation_id = $1`

	tag, err := db.Pool.Exec(ctx, query, id, awsDetail)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *reservationDao) UpdateReservationIDForAWS(ctx context.Context, id int64, awsReservationId string) error {
	query := `UPDATE aws_reservation_details SET aws_reservation_id = $2 WHERE reservation_id = $1`

	tag, err := db.Pool.Exec(ctx, query, id, awsReservationId)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *reservationDao) UpdateOperationNameForGCP(ctx context.Context, id int64, gcpOperationName string) error {
	query := `UPDATE gcp_reservation_details SET gcp_operation_name = $2 WHERE reservation_id = $1`

	tag, err := db.Pool.Exec(ctx, query, id, gcpOperationName)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *reservationDao) FinishWithSuccess(ctx context.Context, id int64) error {
	query := `UPDATE reservations SET success = true, finished_at = now() WHERE id = $1`

	tag, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *reservationDao) FinishWithError(ctx context.Context, id int64, errorString string) error {
	query := `UPDATE reservations SET success = false, error = $2, finished_at = now() WHERE id = $1`

	tag, err := db.Pool.Exec(ctx, query, id, errorString)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}

func (x *reservationDao) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM reservations WHERE id = $1`

	tag, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}
	return nil
}
