package db

import (
	"errors"

	"github.com/jackc/pgconn"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html
type PostgresErrorCode string

const (
	UniqueConstraintErrorCode PostgresErrorCode = "23505"
)

func IsPostgresError(err error, code PostgresErrorCode) error {
	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) && PostgresErrorCode(pgErr.Code) == code {
		return pgErr
	}
	return nil
}
