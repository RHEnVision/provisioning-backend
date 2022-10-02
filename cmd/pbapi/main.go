package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/random"

	// HTTP client implementations
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	"github.com/RHEnVision/provisioning-backend/internal/config/parser"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/RHEnVision/provisioning-backend/internal/version"

	// Job queue implementation
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue/dejq"

	// DAO implementation, must be initialized before any database packages.
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/pgx"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/cloudwatchlogs"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	m "github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func init() {
	random.SeedGlobal()
}

func statusOk(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	ctx := context.Background()
	config.Initialize()

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

	// report unknown environmental variables (see ./internal/config/parser/known.go)
	unknown := parser.UnknownEnvVariables()
	if len(unknown) > 0 {
		logger.Warn().Msgf("Unknown ENV variables, add them in the codebase: %+v", unknown)
	}

	tel := telemetry.Initialize(&log.Logger)
	defer tel.Close(ctx)

	// initialize the rest
	err = db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	// initialize the job queue
	err = dejq.Initialize(ctx, &logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing dejq queue")
	}
	if config.Worker.Queue == "memory" {
		dejq.RegisterJobs(&logger)
	}
	dejq.StartDequeueLoop(ctx, &logger)

	// Routes for the main service
	r := chi.NewRouter()
	r.Use(m.NewPatternMiddleware(version.PrometheusLabelName))
	r.Use(telemetry.Middleware(r))
	r.Use(m.VersionMiddleware)
	r.Use(m.TraceID)
	r.Use(m.LoggerMiddleware(&log.Logger))

	// Set Content-Type to JSON for chi renderer. Warning: Non-chi routes
	// MUST set Content-Type header on their own!
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Setup optional compressor, chi uses sync.Pool so this is cheap.
	// This setup only uses the default gzip which is widely supported
	// across the globe, including HTTP proxies which do have problems with
	// modern algorithms like brotli or zstd to this day.
	// This middleware must be inserted after Content-Type header is set.
	if config.Application.Compression {
		compressor := middleware.NewCompressor(5,
			"application/json",
			"application/x-yaml",
			"text/plain")
		r.Use(compressor.Handler)
	}

	// Unauthenticated routes
	r.Get("/", statusOk)
	// Main routes
	routes.SetupRoutes(r)

	// Routes for metrics
	mr := chi.NewRouter()
	mr.Get("/", statusOk)
	mr.Handle(config.Prometheus.Path, promhttp.Handler())

	log.Info().Msgf("Starting new instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Application.Port),
		Handler: r,
	}

	msrv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Prometheus.Port),
		Handler: mr,
	}

	waitForSignal := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Main service shutdown error")
		}
		if err := msrv.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Metrics service shutdown error")
		}
		close(waitForSignal)
	}()

	go func() {
		if err := msrv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("Metrics service listen error")
			}
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Main service listen error")
		}
	}

	<-waitForSignal

	if config.Worker.Queue == "memory" {
		dejq.StopDequeueLoop()
	}
	log.Info().Msg("Shutdown finished, exiting")
}
