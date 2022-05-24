package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/jmoiron/sqlx"
)

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn` or when it panics.
func WithTransaction(ctx context.Context, db *sqlx.DB, fn TxFn) error {
	logger := ctxval.GetLogger(ctx)
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Cannot begin database transaction")
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			logger.Trace().Msg("Rolling database transaction back")
			err := tx.Rollback()
			if err != nil {
				logger.Error().Err(err).Msg("Cannot rollback database transaction")
				return
			}
			panic(p)
		} else if err != nil {
			logger.Trace().Msg("Rolling database transaction back")
			err := tx.Rollback()
			if err != nil {
				logger.Error().Err(err).Msg("Cannot rollback database transaction")
				return
			}
		} else {
			logger.Trace().Msg("Committing database transaction")
			err = tx.Commit()
			if err != nil {
				logger.Error().Err(err).Msg("Cannot rollback database transaction")
				return
			}
		}
	}()

	logger.Trace().Msg("Starting database transaction")
	err = fn(tx)
	return err
}
