package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/cloudwatchlogs"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func processMessage(ctx context.Context, message *kafka.GenericMessage) {
	// TODO: implement the source checking for AWS, Azure, GCP

	// This needs to be done in 4 goroutines.
	// 1) Reads messages, calls sources to get auth type and ARN and then enqueue the work
	// in one of three channels (AWS, Azure, GCP)
	// 2) Reads from AWS channel, performs availability check, has configurable throttling.
	// 3) Reads from Azure channel, performs availability check, has configurable throttling.
	// 4) Reads from GCP channel, performs availability check, has configurable throttling.

	// Beware there is no database connection, no job queue, most of what api/worker processes have
	// is not available because it should not be needed. The worker will simply use SDKs to check
	// statuses and send out results via kafka.Send function to Sources.
	// Also make sure that all goroutines are closed correctly on context cancel. A WaitGroup
	// would be a really good here: https://gist.github.com/fracasula/b579d52daf15426e58aa133d0340ccb0
}

func main() {
	ctx := context.Background()
	config.Initialize("config/api.env", "config/statuser.env")

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

	// setup structured logging
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}
	logger = logger.With().
		Timestamp().
		Str("hostname", hostname).
		Bool("statuser", true).
		Logger()

	// initialize telemetry
	tel := telemetry.Initialize(&log.Logger)
	defer tel.Close(ctx)

	// initialize platform kafka
	logger.Info().Msg("Initializing platform kafka")
	err = kafka.InitializeKafkaBroker()
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to initialize the platform kafka")
	}

	// start the consumer
	cancelCtx, consumerCancelFunc := context.WithCancel(ctx)
	go kafka.Consume(cancelCtx, kafka.AvailabilityStatusRequestTopic, processMessage)

	// metrics
	logger.Info().Msgf("Starting new instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle(config.Prometheus.Path, promhttp.Handler())
	metricsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Prometheus.Port),
		Handler: metricsRouter,
	}

	waitForSignal := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if err := metricsServer.Shutdown(context.Background()); err != nil {
			logger.Fatal().Err(err).Msg("Metrics service shutdown error")
		}
		close(waitForSignal)
	}()

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			var errInUse syscall.Errno
			if errors.As(err, &errInUse) && errInUse == syscall.EADDRINUSE {
				logger.Warn().Err(err).Msg("Not starting metrics service, port already in use")
			} else if !errors.Is(err, http.ErrServerClosed) {
				logger.Warn().Err(err).Msg("Metrics service listen error")
			}
		}
	}()

	logger.Info().Msg("Worker started")
	<-waitForSignal
	logger.Info().Msg("Consumer shutdown initiated")
	consumerCancelFunc()
	logger.Info().Msg("Shutdown finished, exiting")
}
