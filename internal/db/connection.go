package db

import (
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	// DB is the main connection pool (sqlx on top of database/sql connection pool)
	DB *sqlx.DB
)

func Initialize() error {
	var err error

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		config.GetLoggingConfig().Database.User,
		config.GetLoggingConfig().Database.Password,
		config.GetLoggingConfig().Database.Host,
		config.GetLoggingConfig().Database.Port,
		config.GetLoggingConfig().Database.Name)

	connConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return errors.Wrap(err, "unable to parse database configuration")
	}
	connConfig.Logger = zerologadapter.NewLogger(log.Logger)
	connConfig.LogLevel = pgx.LogLevel(config.GetLoggingConfig().Database.LogLevel)

	DB, err = sqlx.Open("pgx", connStr)
	if err != nil {
		return errors.Wrap(err, "unable to connect to database")
	}

	DB.SetMaxIdleConns(config.GetLoggingConfig().Database.MaxIdleConn)
	DB.SetMaxOpenConns(config.GetLoggingConfig().Database.MaxOpenConn)
	DB.SetConnMaxLifetime(config.GetLoggingConfig().Database.MaxLifetime)
	DB.SetConnMaxIdleTime(config.GetLoggingConfig().Database.MaxIdleTime)
	err = DB.Ping()
	if err != nil {
		return errors.Wrap(err, "unable to ping the database")
	}

	return nil
}
