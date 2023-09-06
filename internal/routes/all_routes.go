package routes

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/api"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/preload"
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	redoc "github.com/go-openapi/runtime/middleware"
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

	// Please make sure this is not gziped, Azure does not like it
	r.Get("/azure_offering_template", s.AzureOfferingTemplate)
	r.Options("/azure_offering_template", s.AzureOfferingTemplate)

	// Review permissions in https://github.com/RedHatInsights/rbac-config when editing this group
	r.Group(func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))

		r.Use(middleware.EnforceIdentity)
		r.Use(middleware.AccountMiddleware)

		// OpenAPI documented and supported routes
		r.Route("/sources", func(r chi.Router) {
			// https://issues.redhat.com/browse/HMS-2305
			// r.Use(middleware.EnforcePermissions("source", "read"))

			r.Get("/", s.ListSources)
			r.With(middleware.Pagination).Get("/", s.ListSources)
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/status", s.SourcesStatus)

				// TODO DEPRECATED: replaced with upload_info
				r.Get("/account_identity", s.GetAWSAccountIdentity)

				r.With(middleware.Pagination).Get("/launch_templates", s.ListLaunchTemplates)
				r.Get("/upload_info", s.GetSourceUploadInfo)
				r.Route("/validate_permissions", func(r chi.Router) {
					r.Get("/", s.ValidatePermissions)
				})
			})
		})

		r.Route("/pubkeys", func(r chi.Router) {
			r.With(middleware.EnforcePermissions("pubkey", "write")).Post("/", s.CreatePubkey)
			r.With(middleware.EnforcePermissions("pubkey", "read")).With(middleware.Pagination).Get("/", s.ListPubkeys)
			r.Post("/", s.CreatePubkey)
			r.Route("/{ID}", func(r chi.Router) {
				r.With(middleware.EnforcePermissions("pubkey", "read")).Get("/", s.GetPubkey)
				r.With(middleware.EnforcePermissions("pubkey", "write")).Delete("/", s.DeletePubkey)
			})
		})

		r.Route("/reservations", func(r chi.Router) {
			r.With(middleware.EnforcePermissions("reservation", "read")).With(middleware.Pagination).Get("/", s.ListReservations)
			// Different types do have different payloads, therefore TYPE must be part of
			// URL and not a URL (filter) parameter.
			r.Route("/{TYPE}", func(r chi.Router) {
				// additional permission checks are in the service functions
				r.With(middleware.EnforcePermissions("reservation", "read")).Get("/{ID}", s.GetReservationDetail)
				r.With(middleware.EnforcePermissions("reservation", "write")).Post("/", s.CreateReservation)
			})
			// Generic reservation detail request (no details provided)
			r.With(middleware.EnforcePermissions("reservation", "read")).Get("/{ID}", s.GetReservationDetail)
		})

		// Endpoint used by sources background checker (no permissions needed)
		r.Route("/availability_status", func(r chi.Router) {
			r.Route("/sources", func(r chi.Router) {
				r.Post("/", s.AvailabilityStatus)
			})
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
				r.Use(middleware.ETagMiddleware(preload.AzureInstanceType.ETagValue))
				r.Get("/", s.ListBuiltinInstanceTypes(preload.AzureInstanceType.InstanceTypesForZone))
			})
			r.Route("/aws", func(r chi.Router) {
				r.Use(middleware.ETagMiddleware(preload.EC2InstanceType.ETagValue))
				r.Get("/", s.ListBuiltinInstanceTypes(preload.EC2InstanceType.InstanceTypesForZone))
			})
			r.Route("/gcp", func(r chi.Router) {
				r.Use(middleware.ETagMiddleware(preload.GCPInstanceType.ETagValue))
				r.Get("/", s.ListBuiltinInstanceTypes(preload.GCPInstanceType.InstanceTypesForZone))
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
