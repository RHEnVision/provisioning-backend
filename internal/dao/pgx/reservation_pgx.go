package pgx

import (
	"context"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

func init() {
	dao.GetReservationDao = getReservationDao
}

type reservationDao struct{}

func getReservationDao(ctx context.Context) dao.ReservationDao {
	return &reservationDao{}
}

func (x *reservationDao) CreateNoop(ctx context.Context, reservation *models.NoopReservation) error {
	reservation.Provider = models.ProviderTypeNoop
	if err := x.createGenericReservation(ctx, &reservation.Reservation); err != nil {
		return fmt.Errorf("failed to create reservation record: %w", err)
	}
	return nil
}

func (x *reservationDao) CreateAWS(ctx context.Context, reservation *models.AWSReservation) error {
	txErr := dao.WithTransaction(ctx, func(tx pgx.Tx) error {
		reservation.Provider = models.ProviderTypeAWS
		if err := x.createGenericReservation(ctx, &reservation.Reservation); err != nil {
			return err
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
		return fmt.Errorf("pgx tx error: %w", txErr)
	}
	return nil
}

func (x *reservationDao) CreateAzure(ctx context.Context, reservation *models.AzureReservation) error {
	txErr := dao.WithTransaction(ctx, func(tx pgx.Tx) error {
		reservation.Provider = models.ProviderTypeAzure
		if err := x.createGenericReservation(ctx, &reservation.Reservation); err != nil {
			return err
		}

		azureQuery := `INSERT INTO azure_reservation_details (reservation_id, pubkey_id, source_id, image_id, detail)
			VALUES ($1, $2, $3, $4, $5)`
		tag, err := db.Pool.Exec(ctx, azureQuery,
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
		return fmt.Errorf("pgx tx error: %w", txErr)
	}
	return nil
}

func (x *reservationDao) CreateGCP(ctx context.Context, reservation *models.GCPReservation) error {
	txErr := dao.WithTransaction(ctx, func(tx pgx.Tx) error {
		reservation.Provider = models.ProviderTypeGCP
		if err := x.createGenericReservation(ctx, &reservation.Reservation); err != nil {
			return err
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
		return fmt.Errorf("pgx tx error: %w", txErr)
	}
	return nil
}

func (x *reservationDao) createGenericReservation(ctx context.Context, reservation *models.Reservation) error {
	reservation.AccountID = identity.AccountId(ctx)
	reservation.Status = "Created"

	reservationQuery := `INSERT INTO reservations (provider, account_id, steps, step_titles, status)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := db.Pool.QueryRow(ctx, reservationQuery,
		reservation.Provider,
		reservation.AccountID,
		reservation.Steps,
		reservation.StepTitles,
		reservation.Status).Scan(&reservation.ID, &reservation.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "too many pending reservations") {
			return fmt.Errorf("%w: %s", dao.ErrReservationRateExceeded, err.Error())
		}
		return fmt.Errorf("failed to create reservation record: %w", err)
	}

	return nil
}

func (x *reservationDao) CreateInstance(ctx context.Context, instance *models.ReservationInstance) error {
	query := `INSERT INTO reservation_instances (reservation_id, instance_id, detail) VALUES ($1, $2, $3)`

	tag, err := db.Pool.Exec(ctx, query,
		instance.ReservationID,
		instance.InstanceID,
		instance.Detail)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row: %w", dao.ErrAffectedMismatch)
	}

	return nil
}

func (x *reservationDao) UpdateReservationInstance(ctx context.Context, reservationID int64, instance *clients.InstanceDescription) error {
	query := `UPDATE reservation_instances SET detail = $3 WHERE reservation_id = $1 AND instance_id = $2`
	detail := &models.ReservationInstanceDetail{
		PublicIPv4: instance.PublicIPv4,
		PublicDNS:  instance.PublicDNS,
	}
	tag, err := db.Pool.Exec(ctx, query, reservationID, instance.ID, detail)
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
	accountId := identity.AccountId(ctx)
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
	accountId := identity.AccountId(ctx)
	result := &models.AWSReservation{}

	err := pgxscan.Get(ctx, db.Pool, result, query, accountId, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) GetAzureById(ctx context.Context, id int64) (*models.AzureReservation, error) {
	query := `SELECT id, reservations.provider, account_id, created_at, steps, step, status, error, finished_at, success,
    	pubkey_id, source_id, image_id, detail
		FROM reservations, azure_reservation_details
		WHERE account_id = $1 AND id = $2 AND id = reservation_id AND reservations.provider = provider_type_azure() LIMIT 1`
	accountId := identity.AccountId(ctx)
	result := &models.AzureReservation{}

	err := pgxscan.Get(ctx, db.Pool, result, query, accountId, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) GetGCPById(ctx context.Context, id int64) (*models.GCPReservation, error) {
	query := `SELECT id, provider, account_id, created_at, steps, step, status, error, finished_at, success,
    	pubkey_id, source_id, image_id, detail
		FROM reservations, gcp_reservation_details
		WHERE account_id = $1 AND id = $2 AND id = reservation_id AND provider = provider_type_gcp() LIMIT 1`
	accountId := identity.AccountId(ctx)
	result := &models.GCPReservation{}

	err := pgxscan.Get(ctx, db.Pool, result, query, accountId, id)
	if err != nil {
		return nil, fmt.Errorf("pgx error: %w", err)
	}
	return result, nil
}

func (x *reservationDao) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM reservations WHERE account_id = $1`
	accountId := identity.AccountId(ctx)

	var result int
	err := db.Pool.QueryRow(ctx, query, accountId).Scan(&result)
	if err != nil {
		return 0, fmt.Errorf("pgx error: %w", err)
	}

	return result, nil
}

func (x *reservationDao) List(ctx context.Context, limit, offset int64) ([]*models.Reservation, error) {
	query := `SELECT * FROM reservations WHERE account_id = $1 ORDER BY id LIMIT $2 OFFSET $3`

	accountId := identity.AccountId(ctx)
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
	query := `SELECT reservation_id, instance_id, detail FROM reservation_instances, reservations
         WHERE reservation_id = reservations.id AND account_id = $1 AND reservation_id = $2`

	accountId := identity.AccountId(ctx)
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

func (x *reservationDao) Cleanup(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	query := `DELETE FROM reservations WHERE created_at < now() - cast($1 as interval)`
	reservationLifetime := config.Reservation.Lifetime.String()

	tag, err := db.Pool.Exec(ctx, query, reservationLifetime)
	if err != nil {
		return fmt.Errorf("pgx error: %w", err)
	}
	logger.Trace().Msgf("Deleted %d reservation(s) older than %s", tag.RowsAffected(), reservationLifetime)

	return nil
}
