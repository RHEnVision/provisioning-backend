package main

import (
	"context"
	"os"

	// DAO implementation, must be initialized before any database packages.
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/migrations"
	"github.com/RHEnVision/provisioning-backend/internal/random"
	"github.com/rs/zerolog/log"
)

func init() {
	random.SeedGlobal()
}

func main() {
	ctx := context.Background()
	config.Initialize("config/api.env", "config/migrate.env")

	// initialize stdout logging and AWS clients first (cloudwatch is not available in init containers)
	logging.InitializeStdout()
	logging.DumpConfigForDevelopment()
	logger := log.Logger

	err := db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	if len(os.Args[1:]) > 0 && os.Args[1] == "purgedb" {
		logger.Warn().Msg("Database purge: all data is being dropped")
		err = migrations.Seed(ctx, "drop_all")
		if err != nil {
			logger.Fatal().Err(err).Msg("Error purging the database")
			return
		}
		logger.Info().Msgf("Database %s has been purged to blank state", config.Database.Name)
	}

	err = migrations.Migrate(ctx, "public")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error running migration")
		return
	}

	if config.Database.SeedScript != "" {
		err = migrations.Seed(ctx, config.Database.SeedScript)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running migration")
			return
		}
	}
}
