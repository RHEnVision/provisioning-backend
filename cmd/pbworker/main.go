package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/RHEnVision/provisioning-backend/internal/random"

	"github.com/RHEnVision/provisioning-backend/internal/config"

	// HTTP client implementations
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"

	// Performs initialization of DAO implementation, must be initialized before any database packages.
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"

	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/zerolog/log"
)

func init() {
	random.SeedGlobal()
}

func main() {
	ctx := context.Background()
	config.Initialize("config/api.env", "config/worker.env")

	// initialize stdout logging and AWS clients first
	logging.InitializeStdout()

	// initialize cloudwatch using the AWS clients
	logger, clsFunc, err := logging.InitializeCloudwatch(log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing cloudwatch")
	}
	defer clsFunc()
	log.Logger = logger
	logging.DumpConfigForDevelopment()

	// setup structured logging
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}
	logger = logger.With().
		Timestamp().
		Str("hostname", hostname).
		Logger()
	logger.Info().Msg("Worker starting")

	// initialize cache
	cache.Initialize()

	// initialize the database
	logger.Debug().Msg("Initializing database connection")
	err = db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	// initialize the job queue
	err = jq.Initialize(ctx, &logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing job queue")
	}
	jq.RegisterJobs(&logger)
	jq.StartDequeueLoop(ctx)

	// wait for term signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info().Msg("Graceful shutdown initiated - waiting for jobs to finish")
	jq.StopDequeueLoop(ctx)
	logger.Info().Msg("Graceful shutdown finished - exiting")
}
