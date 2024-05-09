package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RHEnVision/provisioning-backend/internal/background"
	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/notifications"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func worker() {
	ctx := context.Background()
	cfgs := []string {
		"config/api.env",
		"config/worker.env",
	}
	if len(os.Args) == 3 {
		cfgs = append(cfgs, os.Args[2])
	}
	config.Initialize(cfgs...)

	// initialize stdout logging and AWS clients first
	logging.InitializeStdout()

	// initialize cloudwatch using the AWS clients
	logger, closeFunc := logging.InitializeLogger()
	defer closeFunc()
	logging.DumpConfigForDevelopment()

	logger.Info().Msg("Worker starting")

	// initialize telemetry
	tel := telemetry.Initialize(ctx, &log.Logger)
	defer tel.Close(ctx)
	metrics.RegisterWorkerMetrics()

	// initialize cache
	cache.Initialize()

	// initialize the database
	logger.Debug().Msg("Initializing database connection")
	err := db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	// initialize platform kafka and notifications
	if config.Kafka.Enabled {
		err = kafka.InitializeKafkaBroker(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to initialize the platform kafka")
		}

		if config.Application.Notifications.Enabled {
			notifications.Initialize(ctx)
		}
	}

	// initialize the job queue
	err = jq.Initialize(ctx, &logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing job queue")
	}
	jq.RegisterJobs(&logger)
	jq.StartDequeueLoop(ctx)

	// initialize background goroutines
	bgCtx, bgCancel := context.WithCancel(ctx)
	background.InitializeWorker(bgCtx)
	defer bgCancel()

	// metrics
	logger.Info().Msgf("Starting new instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle(config.Prometheus.Path, promhttp.Handler())
	metricsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Prometheus.Port),
		Handler: metricsRouter,
	}

	signalNotify := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		if err := metricsServer.Shutdown(context.Background()); err != nil {
			logger.Warn().Err(err).Msg("Metrics service shutdown error")
		}
		close(signalNotify)
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

	// wait for term signal
	<-signalNotify

	logger.Info().Msg("Graceful shutdown initiated - waiting for jobs to finish")
	jq.StopDequeueLoop(ctx)
	logger.Info().Msg("Graceful shutdown finished - exiting")
}
