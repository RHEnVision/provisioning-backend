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

	r := chi.NewRouter()
	r.Use(m.RequestID)
	r.Use(m.RequestNum)
	r.Use(m.MetricsMiddleware)
	r.Use(m.LoggerMiddleware(&log.Logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Handle("/metrics", promhttp.Handler())
	routes.SetupRoutes(r)

	log.Info().Msg("New instance started ")
	http.ListenAndServe(":3000", r)
}
