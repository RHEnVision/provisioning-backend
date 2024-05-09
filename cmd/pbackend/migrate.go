package main

import (
	"context"
	"os"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/migrations"
	"github.com/rs/zerolog/log"
)

func migrate() {
	ctx := context.Background()
	cfgs := []string {
		"config/api.env",
		"config/migrate.env",
	}
	if len(os.Args) == 3 {
		cfgs = append(cfgs, os.Args[2])
	}
	config.Initialize(cfgs...)

	// initialize stdout logging and AWS clients first (cloudwatch is not available in init containers)
	logging.InitializeStdout()
	logging.DumpConfigForDevelopment()
	logger := log.Logger

	err := db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	if len(os.Args[2:]) > 0 && os.Args[2] == "purgedb" {
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
