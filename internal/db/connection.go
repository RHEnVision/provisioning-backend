package db

import (
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
)

var (
	// DB is the main connection pool (sqlx on top of database/sql connection pool)
	DB *sqlx.DB
)

func GetConnectionString(prefix string) string {
	if len(config.Database.Password) > 0 {
		return fmt.Sprintf("%s://%s:%s@%s:%d/%s",
			prefix,
			url.QueryEscape(config.Database.User),
			url.QueryEscape(config.Database.Password),
			config.Database.Host,
			config.Database.Port,
			config.Database.Name)
	} else {
		return fmt.Sprintf("%s://%s@%s:%d/%s",
			prefix,
			url.QueryEscape(config.Database.User),
			config.Database.Host,
			config.Database.Port,
			config.Database.Name)
	}

}
func Initialize() error {
	var err error

	connStr := GetConnectionString("postgres")
	connConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return errors.Wrap(err, "unable to parse database configuration")
	}
	connConfig.Logger = zerologadapter.NewLogger(log.Logger)
	connConfig.LogLevel = pgx.LogLevel(config.Database.LogLevel)

	DB, err = sqlx.Open("pgx", connStr)
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
