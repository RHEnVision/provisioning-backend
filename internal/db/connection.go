package db

import (
	"fmt"
	"net/url"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	// DB is the main connection pool (sqlx on top of database/sql connection pool)
	DB *sqlx.DB
)

func GetConnectionString(prefix, schema string) string {
	if len(config.Database.Password) > 0 {
		return fmt.Sprintf("%s://%s:%s@%s:%d/%s?search_path=%s",
			prefix,
			url.QueryEscape(config.Database.User),
			url.QueryEscape(config.Database.Password),
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			schema)
	} else {
		return fmt.Sprintf("%s://%s@%s:%d/%s?search_path=%s",
			prefix,
			url.QueryEscape(config.Database.User),
			config.Database.Host,
			config.Database.Port,
			config.Database.Name,
			schema)
	}

}
func Initialize(schema string) error {
	var err error
	if schema == "" {
		schema = "public"
	}

	// register and setup logging configuration
	connStr := GetConnectionString("postgres", schema)
	connConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return errors.Wrap(err, "unable to parse database configuration")
	}
	if config.Database.LogLevel > 0 {
		connConfig.Logger = zerologadapter.NewLogger(log.Logger)
		connConfig.LogLevel = pgx.LogLevel(config.Database.LogLevel)
	}
	connStrRegistered := stdlib.RegisterConnConfig(connConfig)

	DB, err = sqlx.Open("pgx", connStrRegistered)
	if err != nil {
		return errors.Wrap(err, "unable to connect to database")
	}

	DB.SetMaxIdleConns(config.Database.MaxIdleConn)
	DB.SetMaxOpenConns(config.Database.MaxOpenConn)
	DB.SetConnMaxLifetime(config.Database.MaxLifetime)
	DB.SetConnMaxIdleTime(config.Database.MaxIdleTime)
	err = DB.Ping()
	if err != nil {
		return errors.Wrap(err, "unable to ping the database")
	}

	return nil
}
