package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/IBM/pgxpoolprometheus"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/exaring/otelpgx"
	pgxlog "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/georgysavva/scany/v2"
)

// Pool is the main connection pool for the whole application
var Pool *pgxpool.Pool

func getConnString(prefix, schema string) string {
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

// Initialize creates connection pool. Close must be called when done.
func Initialize(ctx context.Context, schema string) error {
	var err error
	if schema == "" {
		schema = "public"
	}

	// register and setup logging configuration
	connStr := getConnString("postgres", schema)
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("unable to parse db configuration: %w", err)
	}

	poolConfig.MaxConns = config.Database.MaxConn
	poolConfig.MinConns = config.Database.MinConn
	poolConfig.MaxConnLifetime = config.Database.MaxLifetime
	poolConfig.MaxConnIdleTime = config.Database.MaxIdleTime

	if config.Telemetry.Enabled {
		poolConfig.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithIncludeQueryParameters())
	} else {
		logLevel, configErr := tracelog.LogLevelFromString(config.Database.LogLevel)
		if configErr != nil {
			return fmt.Errorf("cannot parse db log level configuration: %w", configErr)
		}

		if logLevel > 0 {
			zeroLogger := pgxlog.NewLogger(log.Logger,
				pgxlog.WithContextFunc(func(ctx context.Context, zx zerolog.Context) zerolog.Context {
					jobId := logging.JobId(ctx)
					if jobId != "" {
						zx = zx.Str("job_id", jobId)
					}
					reservationId := logging.ReservationId(ctx)
					if reservationId != 0 {
						zx = zx.Int64("reservation_id", reservationId)
					}
					traceId := logging.TraceId(ctx)
					if traceId != "" {
						zx = zx.Str("trace_id", traceId)
					}
					requestId := logging.EdgeRequestId(ctx)
					if requestId != "" {
						zx = zx.Str("request_id", requestId)
					}
					accountId := identity.AccountIdOrNil(ctx)
					if accountId != 0 {
						zx = zx.Int64("account_id", accountId)
					}
					principal := identity.Identity(ctx)
					zx = zx.Str("org_id", principal.Identity.OrgID)
					zx = zx.Str("account_number", principal.Identity.AccountNumber)
					return zx
				}))
			poolConfig.ConnConfig.Tracer = &tracelog.TraceLog{
				Logger:   zeroLogger,
				LogLevel: logLevel,
			}
		}
	}

	Pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = Pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("unable to ping the database: %w", err)
	}

	// Register telemetry
	labels := map[string]string{
		"service": version.PrometheusLabelName,
		"db_host": config.Database.Host,
		"db_name": config.Database.Name,
	}
	collector := pgxpoolprometheus.NewCollector(Pool, labels)
	prometheus.MustRegister(collector)

	return nil
}

func Close() {
	log.Logger.Info().Msg("Closing all database connections")
	Pool.Close()
}
