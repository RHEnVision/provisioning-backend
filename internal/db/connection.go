package db

import (
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/jobqueue/dbjobqueue"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

var (
	// DB is the main connection pool (sqlx on top of database/sql)
	DB *sqlx.DB

	// JQ is the main job queue
	JQ *dbjobqueue.DBJobQueue
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
	DB.SetConnMaxLifetime(time.Duration(config.GetLoggingConfig().Database.MaxLifetime) * time.Second)
	DB.SetConnMaxIdleTime(time.Duration(config.GetLoggingConfig().Database.MaxIdleTime) * time.Second)
	err = DB.Ping()
	if err != nil {
		return errors.Wrap(err, "unable to ping the database")
	}

	// create new job queue (creates a separate connection pool)
	JQ, err = dbjobqueue.New(connStr)
	if err != nil {
		return errors.Wrap(err, "unable to create job queue")
	}

	return nil
}
