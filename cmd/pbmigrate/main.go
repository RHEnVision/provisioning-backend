package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/clients/cloudwatchlogs"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/zerolog/log"
)

func main() {
	// initialize stdout logging and AWS clients first
	log.Logger = logging.InitializeStdout()
	cloudwatchlogs.Initialize()

	// initialize cloudwatch using the AWS clients
	logger, clsFunc, err := logging.InitializeCloudwatch(log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing cloudwatch")
	}
	defer clsFunc()
	log.Logger = logger
	logging.DumpConfigForDevelopment()

	err = db.Initialize()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}

	log.Info().Msg("Migrating database")
	db.Migrate()
	log.Info().Msgf("Migration complete")
	if config.Database.SeedScript != "" {
		log.Info().Msgf("Seeding '%s'", config.Database.SeedScript)
		db.Seed(config.Database.SeedScript)
		log.Info().Msg("Seed complete")
	}
}
