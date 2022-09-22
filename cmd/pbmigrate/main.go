package main

import (
	"os"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/cloudwatchlogs"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/random"
	"github.com/rs/zerolog/log"
)

func init() {
	random.SeedGlobal()
}

func main() {
	config.Initialize()

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

	err = db.Initialize("public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}

	if len(os.Args[1:]) > 0 && os.Args[1] == "purgedb" {
		logger.Warn().Msg("Database purge: all data is being dropped")
		err = db.Seed("drop_all")
		if err != nil {
			logger.Fatal().Err(err).Msg("Error purging the database")
			return
		}
		logger.Info().Msgf("Database %s has been purged to blank state", config.Database.Name)
	}

	err = db.Migrate("public")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error running migration")
		return
	}

	if config.Database.SeedScript != "" {
		err = db.Seed(config.Database.SeedScript)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running migration")
			return
		}
	}
}
