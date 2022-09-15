package db

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/tern/migrate"
)

//go:embed migrations
var embeddedMigrations embed.FS

//go:embed seeds
var embeddedSeeds embed.FS

type EmbeddedFS struct {
	efs *embed.FS
}

func NewEmbeddedFS(fs *embed.FS) *EmbeddedFS {
	return &EmbeddedFS{efs: fs}
}

func (efs *EmbeddedFS) ReadDir(dirname string) ([]fs.FileInfo, error) {
	dirEntries, err := efs.efs.ReadDir(dirname)
	if err != nil {
		return nil, fmt.Errorf("unable to read dir: %w", err)
	}
	result := make([]fs.FileInfo, 0, len(dirEntries))
	for _, de := range dirEntries {
		fi, err := de.Info()
		if err != nil {
			return nil, fmt.Errorf("unable to read dir: %w", err)
		}
		result = append(result, fi)
	}
	return result, nil
}

func (efs *EmbeddedFS) ReadFile(filename string) ([]byte, error) {
	result, err := efs.efs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %w", err)
	}
	return result, nil
}

func (efs *EmbeddedFS) Glob(pattern string) (matches []string, err error) {
	result, err := fs.Glob(efs.efs, pattern)
	if err != nil {
		return nil, fmt.Errorf("unable to glob: %w", err)
	}
	return result, nil
}

func fmtDetailedError(sql string, mgErr *pgconn.PgError) string {
	var errb strings.Builder
	errb.WriteString(mgErr.Error())

	if mgErr.Detail != "" {
		errb.WriteString(fmt.Sprintln("DETAIL:", mgErr.Detail))
	}

	if mgErr.Position != 0 {
		ele, err := ExtractErrorLine(sql, int(mgErr.Position))
		if err != nil {
			errb.WriteString(err.Error())
			return errb.String()
		}

		prefix := fmt.Sprintf("\nLINE %d: ", ele.LineNum)
		errb.WriteString(fmt.Sprintf("%s%s\n", prefix, ele.Text))

		padding := strings.Repeat(" ", len(prefix)+ele.ColumnNum-1)
		errb.WriteString(fmt.Sprintf("%s^\n", padding))
	}

	if mgErr.Where != "" {
		errb.WriteString(fmt.Sprintf(", WHERE: %s\n", mgErr.Where))
	}

	if mgErr.InternalPosition != 0 {
		ele, err := ExtractErrorLine(mgErr.InternalQuery, int(mgErr.InternalPosition))
		if err != nil {
			errb.WriteString(err.Error())
			return errb.String()
		}

		prefix := fmt.Sprintf("LINE %d: ", ele.LineNum)
		errb.WriteString(fmt.Sprintf("%s%s\n", prefix, ele.Text))

		padding := strings.Repeat(" ", len(prefix)+ele.ColumnNum-1)
		errb.WriteString(fmt.Sprintf("%s^\n", padding))
	}

	return errb.String()
}

var (
	ErrNoMigrationsFound = errors.New("no migrations found")
	ErrMigration         = errors.New("unable to perform migration")
	ErrSeedProduction    = errors.New("seed in production")
)

// Migrate executes embedded SQL scripts from internal/db/migrations. For the time being
// only "up" migrations are supported. When this package is initialized, the directory
// is verified that it only contains XXX_*.up.sql files (XXX = numbers).
func Migrate(schema string) error {
	logger := log.Logger.With().Bool("migration", true).Logger()
	ctx := context.Background()
	logger.Debug().Msgf("Started migration")
	if schema == "" {
		schema = "public"
	}

	stdConn, connErr := DB.Conn(ctx)
	if connErr != nil {
		return fmt.Errorf("error acquiring connection from pool: %w", connErr)
	}
	defer stdConn.Close()

	wrapErr := stdConn.Raw(func(pgxConn interface{}) error {
		conn := pgxConn.(*stdlib.Conn).Conn()
		opts := migrate.MigratorOptions{
			MigratorFS: NewEmbeddedFS(&embeddedMigrations),
		}
		table := fmt.Sprintf("%s.schema_version", schema)
		migrator, err := migrate.NewMigratorEx(ctx, conn, table, &opts)
		if err != nil {
			return fmt.Errorf("error initializing migrator: %w", err)
		}
		err = migrator.LoadMigrations("migrations/")
		if err != nil {
			return fmt.Errorf("error loading migrations: %w", err)
		}
		if len(migrator.Migrations) == 0 {
			return ErrNoMigrationsFound
		}

		migrator.OnStart = func(sequence int32, name, direction, sql string) {
			logger.Info().Str("sql", sql).Msgf("Executing migration %s %s", name, direction)
		}

		err = migrator.Migrate(ctx)
		if err != nil {
			var mgErr *migrate.MigrationPgError
			var pgErr *pgconn.PgError
			if errors.As(err, &mgErr) && errors.As(err, &pgErr) {
				return fmt.Errorf("%w: %s", ErrMigration, fmtDetailedError(mgErr.Sql, pgErr))
			} else {
				return fmt.Errorf("unable to perform migration: %w", err)
			}
		}

		return nil
	})
	if wrapErr != nil {
		return fmt.Errorf("error migrating: %w", wrapErr)
	}

	// Print some additional info
	rows, err := DB.Query("SELECT version, applied_at FROM schema_migrations_history")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error querying schema history")
	}
	defer rows.Close()
	for rows.Next() {
		var version int
		var appliedAt time.Time

		if err := rows.Scan(&version, &appliedAt); err != nil {
			logger.Fatal().Err(err).Msg("Error scanning schema history")
		}
		logger.Info().Msgf("Version %d was applied %v", version, appliedAt.UTC())
	}
	if err := rows.Err(); err != nil {
		logger.Fatal().Err(err).Msg("Error scanning schema history")
	}

	logger.Info().Msgf("Finished with migration")
	return nil
}

// Seed executes embedded SQL scripts from internal/db/seeds
func Seed(seedScript string) error {
	logger := log.Logger.With().Bool("seed", true).Logger()
	logger.Debug().Msgf("Started execution of seed script %s", seedScript)

	// Prevent from accidental execution of drop_all seed in production
	if seedScript == "drop_all" && config.Features.Environment != "development" {
		return fmt.Errorf("%w: an attempt to run drop_all seed script in non %s environment", ErrSeedProduction, config.Features.Environment)
	}
	file, err := embeddedSeeds.Open(fmt.Sprintf("seeds/%s.sql", seedScript))
	if err != nil {
		return fmt.Errorf("unable to open seed script %s: %w", seedScript, err)
	}
	defer file.Close()
	buffer, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read seed script %s: %w", seedScript, err)
	}
	_, err = DB.Exec(string(buffer))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			detail := fmtDetailedError(string(buffer), pgErr)
			logger.Fatal().Err(pgErr).Msg(detail)
		} else {
			return fmt.Errorf("unable to execute script %s: %w", seedScript, err)
		}
	}

	logger.Info().Msgf("Executed seed script %s", seedScript)
	return nil
}
