package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/RHEnVision/provisioning-backend/internal/clouds/aws"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	m "github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/routes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

func statusOk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

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

	// Routes for the main service
	r := chi.NewRouter()
	r.Use(m.RequestID)
	r.Use(m.RequestNum)
	r.Use(m.MetricsMiddleware)
	r.Use(m.LoggerMiddleware(&log.Logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// Unauthenticated routes
	r.Get("/", statusOk)
	// Main routes
	routes.SetupRoutes(r)

	// Routes for metrics
	mr := chi.NewRouter()
	mr.Get("/", statusOk)
	mr.Handle("/metrics", promhttp.Handler())

	log.Info().Msg("Starting new instance on 8000/8080")
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
		Handler: r,
	}

	msrv := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: mr,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Main service shutdown error")
		}
		if err := msrv.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Metrics service shutdown error")
		}
		close(idleConnsClosed)
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

	<-idleConnsClosed
	log.Info().Msg("Shutdown finished, exiting")
}
