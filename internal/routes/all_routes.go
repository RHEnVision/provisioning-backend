package routes

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/api"
	azure_types "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure/types"
	ec2_types "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2/types"
	gcp_types "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp/types"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/go-chi/chi/v5"
	redoc "github.com/go-openapi/runtime/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog/log"
)

func redocMiddleware(handler http.Handler) http.Handler {
	opt := redoc.RedocOpts{
		Title:   "Provisioning OpenAPI",
		SpecURL: fmt.Sprintf("%s/openapi.json", PathPrefix()),
	}
	return redoc.Redoc(opt, handler)
}

func logETags() {
	for _, etag := range middleware.AllETags() {
		log.Logger.Trace().Msgf("Calculated '%s' etag '%s' in %dms", etag.Name, etag.Value, etag.HashTime.Milliseconds())
	}
}

func MountRoot(r *chi.Mux) {
	logETags()

	r.Get("/", s.WelcomeService)
	r.Get("/ping", s.StatusService)
	r.Route("/docs", func(r chi.Router) {
		r.Use(redocMiddleware)
		r.Route("/openapi.json", func(r chi.Router) {
			r.Use(middleware.ETagMiddleware(api.ETagValue))
			r.Get("/", api.ServeOpenAPISpec)
		})
	})
}

func MountAPI(r *chi.Mux) {
	r.Route("/openapi.json", func(r chi.Router) {
		r.Use(middleware.ETagMiddleware(api.ETagValue))
		r.Get("/", api.ServeOpenAPISpec)
	})
	r.Group(func(r chi.Router) {
		r.Use(identity.EnforceIdentity)
		r.Use(middleware.AccountMiddleware)

		// OpenAPI documented and supported routes
		r.Route("/sources", func(r chi.Router) {
			r.Get("/", s.ListSources)
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/status", s.SourcesStatus)

				// TODO move this to outside of /sources (see below)
				r.Get("/instance_types", s.ListInstanceTypes)
			})
		})

		r.Route("/pubkeys", func(r chi.Router) {
			r.Post("/", s.CreatePubkey)
			r.Get("/", s.ListPubkeys)
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", s.GetPubkey)
				r.Delete("/", s.DeletePubkey)
			})
		})

		r.Route("/reservations", func(r chi.Router) {
			r.Get("/", s.ListReservations)
			// Different types do have different payloads, therefore TYPE must be part of
			// URL and not a URL (filter) parameter.
			r.Route("/{TYPE}", func(r chi.Router) {
				r.Get("/{ID}", s.GetReservationDetail)
				r.Post("/", s.CreateReservation)
			})
			// Generic reservation detail request (no details provided)
			r.Get("/{ID}", s.GetReservationDetail)
		})

		// Unsupported routes are not published through OpenAPI, they are documented
		// here. These can be either work-in-progress features, infrastructure or
		// development related.

		// Readiness of the service.
		r.Route("/ready", func(r chi.Router) {
			// Returns immediately, no database connection is made
			r.Get("/", s.ReadyService)

			// Connects to a remote service via HTTP client.
			r.Route("/{SRV}", func(r chi.Router) {
				r.Get("/", s.ReadyBackendService)
			})
		})

		// All embedded instance types which are compiled in the application. Allows
		// filtering by provider, region and availability zone. Uses ETag caching for
		// improved UI experience.
		r.Route("/instance_types", func(r chi.Router) {
			r.Route("/azure", func(r chi.Router) {
				r.Use(middleware.ETagMiddleware(azure_types.ETagValue))
				r.Get("/", s.ListBuiltinInstanceTypes(azure_types.InstanceTypesForZone))
			})
			r.Route("/aws", func(r chi.Router) {
				r.Use(middleware.ETagMiddleware(ec2_types.ETagValue))
				r.Get("/", s.ListBuiltinInstanceTypes(ec2_types.InstanceTypesForZone))
			})
			r.Route("/gcp", func(r chi.Router) {
				r.Use(middleware.ETagMiddleware(gcp_types.ETagValue))
				r.Get("/", s.ListBuiltinInstanceTypes(gcp_types.InstanceTypesForZone))
			})
		})

		// We expose feature flags for image builder, this is undocumented since we
		// want to push for the setup where we share the same unleash instance and this
		// endpoint might not be needed anymore.
		r.Route("/feature/{FLAG}", func(r chi.Router) {
			r.Get("/", s.FeatureFlagService)
			r.Head("/", s.FeatureFlagService)
		})
	})
}
