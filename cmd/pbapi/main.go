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
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/RHEnVision/provisioning-backend/internal/random"
	s "github.com/RHEnVision/provisioning-backend/internal/services"

	// HTTP client implementations
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/RHEnVision/provisioning-backend/internal/version"

	// DAO implementation, must be initialized before any database packages.
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"

	// set GOMAXPROC to ideal value for Kubernetes
	_ "go.uber.org/automaxprocs"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	m "github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func init() {
	random.SeedGlobal()
}

func main() {
	ctx := context.Background()
	config.Initialize("config/api.env")

	// initialize cloudwatch using the AWS clients
	logger, closeFunc := logging.InitializeLogger()
	defer closeFunc()
	log.Logger = logger
	logging.DumpConfigForDevelopment()

	// initialize feature flags
	err := config.InitializeFeatureFlags(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing feature flags")
	}
	defer config.StopFeatureFlags(ctx)

	// initialize telemetry
	tel := telemetry.Initialize(&log.Logger)
	defer tel.Close(ctx)
	metrics.RegisterApiMetrics()

	// initialize the rest
	err = db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	// initialize cache
	cache.Initialize()

	// initialize platform kafka
	if config.Kafka.Enabled {
		err = kafka.InitializeKafkaBroker(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to initialize the platform kafka")
		}
	}

	// initialize background goroutines
	bgCtx, bgCancel := context.WithCancel(ctx)
	background.InitializeApi(bgCtx)
	defer bgCancel()

	// initialize job queue
	err = jq.Initialize(ctx, &logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing job queue")
	}
	jq.RegisterJobs(&logger)
	if config.Worker.Queue == "memory" {
		jq.StartDequeueLoop(ctx)
	}

	// Setup routes
	rootRouter := chi.NewRouter()
	apiRouter := chi.NewRouter()

	apiRouter.Use(m.NewPatternMiddleware(version.PrometheusLabelName))
	apiRouter.Use(telemetry.Middleware(apiRouter))
	apiRouter.Use(m.VersionMiddleware)
	apiRouter.Use(m.CorrelationID)
	apiRouter.Use(m.TraceID)
	apiRouter.Use(m.LoggerMiddleware(&log.Logger))

	// Mount paths
	routes.MountRoot(rootRouter)
	routes.MountAPI(apiRouter)
	rootRouter.Mount(routes.PathPrefix(), apiRouter)

	// Routes for metrics
	metricsRouter := chi.NewRouter()
	metricsRouter.Get("/", s.WelcomeService)
	metricsRouter.Handle(config.Prometheus.Path, promhttp.Handler())

	log.Info().Msgf("Starting new instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	apiServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Application.Port),
		Handler: rootRouter,
	}

	metricsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Prometheus.Port),
		Handler: metricsRouter,
	}

	waitForSignal := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if err := apiServer.Shutdown(context.Background()); err != nil {
			log.Warn().Err(err).Msg("Main service shutdown error")
		}
		if err := metricsServer.Shutdown(context.Background()); err != nil {
			log.Warn().Err(err).Msg("Metrics service shutdown error")
		}
		close(waitForSignal)
	}()

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			var errInUse syscall.Errno
			if errors.As(err, &errInUse) && errInUse == syscall.EADDRINUSE {
				log.Warn().Err(err).Msg("Not starting metrics service, port already in use")
			} else if !errors.Is(err, http.ErrServerClosed) {
				log.Warn().Err(err).Msg("Metrics service listen error")
			}
		}
	}()

	if err := apiServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Warn().Err(err).Msg("Main service listen error")
		}
	}

	<-waitForSignal

	if config.Worker.Queue == "memory" {
		jq.StopDequeueLoop(ctx)
	}
	log.Info().Msg("Shutdown finished, exiting")
}
