package main

import (
	"github.com/RHEnVision/provisioning-backend/internal/clouds/aws"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	m "github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/routes"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

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

	log.Info().Msg("New instance started")
	http.ListenAndServe(":3000", r)
	http.ListenAndServe(":5000", mr)
}
