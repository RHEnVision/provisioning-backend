package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/clouds/aws"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/zerolog/log"
)

func main() {
	// initialize stdout logging and AWS clients first
	log.Logger = logging.InitializeStdout()
	aws.Initialize()

	// initialize cloudwatch using the AWS clients
	logger, clsFunc, err := logging.InitializeCloudwatch(log.Logger)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer clsFunc()
	log.Logger = logger

	// initialize the rest
	db.Initialize()

	log.Info().Msg("Migrating database")
	db.Migrate()
	log.Info().Msg("Migration complete")
}
