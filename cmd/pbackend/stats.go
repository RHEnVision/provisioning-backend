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
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func stats() {
	ctx := context.Background()
	cfgs := []string {
		"config/api.env",
		"config/stats.env",
	}
	if len(os.Args) == 3 {
		cfgs = append(cfgs, os.Args[2])
	}
	config.Initialize(cfgs...)

	// initialize cloudwatch using the AWS clients
	logger, closeFunc := logging.InitializeLogger()
	defer closeFunc()
	logging.DumpConfigForDevelopment()

	// initialize telemetry
	tel := telemetry.Initialize(ctx, &log.Logger)
	defer tel.Close(ctx)

	// initialize the job queue but don't register any workers
	err := jq.Initialize(ctx, &logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing job queue")
	}

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
		if shutdownErr := metricsServer.Shutdown(context.Background()); shutdownErr != nil {
			logger.Warn().Err(shutdownErr).Msg("Metrics service shutdown error")
		}
		close(signalNotify)
	}()

	go func() {
		if listenErr := metricsServer.ListenAndServe(); listenErr != nil {
			var errInUse syscall.Errno
			if errors.As(listenErr, &errInUse) && errInUse == syscall.EADDRINUSE {
				logger.Warn().Err(listenErr).Msg("Not starting metrics service, port already in use")
			} else if !errors.Is(listenErr, http.ErrServerClosed) {
				logger.Warn().Err(listenErr).Msg("Metrics service listen error")
			}
		}
	}()

	metrics.RegisterStatsMetrics()

	// initialize the database
	logger.Debug().Msg("Initializing database connection")
	err = db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	// initialize background goroutines
	bgCtx, bgCancel := context.WithCancel(ctx)
	background.InitializeStats(bgCtx)
	defer bgCancel()

	logger.Info().Msg("Stats process started")
	select {
	case <-signalNotify:
		logger.Info().Msg("Exiting due to signal")
	}

	logger.Info().Msg("Stats process shutdown initiated")
}
