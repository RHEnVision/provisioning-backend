package db

import (
	"embed"
	"log"
	"strconv"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations
var fs embed.FS

func Migrate() {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		log.Fatal(err)
	}
	// TODO: extract configuration
	m, err := migrate.NewWithSourceInstance("iofs", d, "postgres://lzap@nuc/pb_dev?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
}

// Checks that migration files are in proper format and index has no gaps or reused numbers.
func init() {
	dir, err := fs.ReadDir("migrations")
	if err != nil {
		log.Fatal("Unable to open migrations embedded directory")
	}
	if len(dir)%2 != 0 {
		log.Fatal("Number of migration files must be even")
	}
	// count migration prefixes
	checks := make([]int, len(dir)/2)
	for _, de := range dir {
		ix, err := strconv.Atoi(de.Name()[:3])
		if err != nil {
			log.Fatalf("Migration %s does not start with an integer?", de.Name())
		}
		if ix-1 > len(checks)-1 {
			log.Fatalf("Is there a gap in migration numbers? Number %d is way too high", ix)
		}
		checks[ix-1]++
	}
	// check expected result
	for i, x := range checks {
		if x != 2 {
			log.Fatalf("There are not exactly two migration files with index %05d, found: %d", i+1, x)
		}
	}
}
