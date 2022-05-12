package db

import (
	"embed"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	stdlog "log"
	"strconv"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//go:embed migrations
var fs embed.FS

// MigrationLogger implements
// https://github.com/golang-migrate/migrate/blob/master/log.go
type MigrationLogger struct {
	logger zerolog.Logger
}

func NewMigrationLogger(logger zerolog.Logger) *MigrationLogger {
	return &MigrationLogger{logger: logger}
}

func (log *MigrationLogger) Printf(format string, v ...interface{}) {
	log.logger.Info().Msgf(format, v...)
}

// Verbose should return true when verbose logging output is wanted
func (log *MigrationLogger) Verbose() bool {
	return true
}

func Migrate() {
	mlog := log.Logger.With().Bool("migration", true).Logger()
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		mlog.Fatal().Err(err).Msg("Error reading migrations")
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, GetConnectionString("pgx"))
	if err != nil {
		mlog.Fatal().Err(err).Msg("Error connecting to database")
	}
	m.Log = NewMigrationLogger(mlog)

	// Perform migration
	if err := m.Up(); errors.Is(err, migrate.ErrNoChange) {
		mlog.Info().Msg("No changes")
	} else {
		mlog.Fatal().Err(err).Msg("Error performing migrations")
	}

	// Print some additional info
	rows, err := DB.Query("SELECT version, applied_at FROM schema_migrations_history")
	if err != nil {
		mlog.Fatal().Err(err).Msg("Error querying schema history")
	}
	defer rows.Close()
	for rows.Next() {
		var version int
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			mlog.Fatal().Err(err).Msg("Error scanning schema history")
		}
		mlog.Info().Msgf("Version %d was applied %v", version, appliedAt)
	}
	if err := rows.Err(); err != nil {
		mlog.Fatal().Err(err).Msg("Error iterating rows")
	}
}

// Checks that migration files are in proper format and index has no gaps or
// reused numbers. Note this runs during package import, so the main logger
// is not yet available. Typically, this fails before it gets into production
// (e.g. during local testing or on CI).
func init() {
	dir, err := fs.ReadDir("migrations")
	if err != nil {
		stdlog.Fatal("Unable to open migrations embedded directory")
	}
	if len(dir)%2 != 0 {
		stdlog.Fatal("Number of migration files must be even")
	}
	// count migration prefixes
	checks := make([]int, len(dir)/2)
	for _, de := range dir {
		ix, err := strconv.Atoi(de.Name()[:3])
		if err != nil {
			stdlog.Fatalf("Migration %s does not start with an integer?", de.Name())
		}
		if ix-1 > len(checks)-1 {
			stdlog.Fatalf("Is there a gap in migration numbers? Number %d is way too high", ix)
		}
		checks[ix-1]++
	}
	// check expected result
	for i, x := range checks {
		if x != 2 {
			stdlog.Fatalf("There are not exactly two migration files with index %05d, found: %d", i+1, x)
		}
	}
}
