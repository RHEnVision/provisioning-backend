package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	// HTTP client implementations
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/image_builder"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/sources"

	// Cache implementation
	_ "github.com/RHEnVision/provisioning-backend/internal/cache/mem"

	// Performs initialization of DAO implementation, must be initialized before any database packages.
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/sqlx"

	// Job queue implementation
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue/dejq"

	"github.com/RHEnVision/provisioning-backend/internal/clients/cloudwatchlogs"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
)

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

	// setup structured logging
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}
	logger = logger.With().
		Timestamp().
		Str("hostname", hostname).
		Str("worker_id", xid.New().String()).
		Logger()
	logger.Info().Msg("Worker starting")

	// initialize the database
	logger.Debug().Msg("Initializing database connection")
	err = db.Initialize("public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}

	// initialize the job queue
	ctx := context.Background()
	err = dejq.Initialize(ctx, &logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing dejq queue")
	}
	dejq.RegisterJobs(&logger)
	dejq.StartDequeueLoop(ctx, &logger)

	// wait for term signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	logger.Info().Msg("Graceful shutdown initiated - waiting for jobs to finish")
	dejq.StopDequeueLoop()
	logger.Info().Msg("Graceful shutdown finished - exiting")
}
