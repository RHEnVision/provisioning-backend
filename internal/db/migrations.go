package db

import (
	"embed"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	stdlog "log"
	"strconv"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
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
	log.logger.Info().Msgf(format, v)
}

// Verbose should return true when verbose logging output is wanted
func (log *MigrationLogger) Verbose() bool {
	return true
}

func Migrate() {
	mlog := log.Logger.With().Bool("migration", true).Logger()
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		mlog.Fatal().Err(err)
	}
	// TODO: extract configuration
	m, err := migrate.NewWithSourceInstance("iofs", d, "postgres://lzap@nuc/pb_dev?sslmode=disable")
	m.Log = NewMigrationLogger(mlog)
	if err != nil {
		mlog.Fatal().Err(err)
	}
	err = m.Up()
	if err != nil {
		mlog.Fatal().Err(err)
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
