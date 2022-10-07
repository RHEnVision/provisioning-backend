package main

import (
	"context"
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
	ctx := context.Background()
	config.Initialize("configs/api.env", "configs/migrate.env")

	// initialize stdout logging and AWS clients first
	logging.InitializeStdout()
	cloudwatchlogs.Initialize()

	// initialize cloudwatch using the AWS clients
	logger, clsFunc, err := logging.InitializeCloudwatch(log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing cloudwatch")
	}
	defer clsFunc()
	log.Logger = logger
	logging.DumpConfigForDevelopment()

	err = db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	if len(os.Args[1:]) > 0 && os.Args[1] == "purgedb" {
		logger.Warn().Msg("Database purge: all data is being dropped")
		err = db.Seed(ctx, "drop_all")
		if err != nil {
			logger.Fatal().Err(err).Msg("Error purging the database")
			return
		}
		logger.Info().Msgf("Database %s has been purged to blank state", config.Database.Name)
	}

	err = db.Migrate(ctx, "public")
	if err != nil {
		logger.Fatal().Err(err).Msg("Error running migration")
		return
	}

	if config.Database.SeedScript != "" {
		err = db.Seed(ctx, config.Database.SeedScript)
		if err != nil {
			logger.Fatal().Err(err).Msg("Error running migration")
			return
		}
	}
}
