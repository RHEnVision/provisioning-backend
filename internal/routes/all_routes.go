package routes

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/api"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/go-chi/chi/v5"
	redoc "github.com/go-openapi/runtime/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func redocMiddleware(handler http.Handler) http.Handler {
	opt := redoc.RedocOpts{
		SpecURL: fmt.Sprintf("%s/openapi.json", PathPrefix()),
	}
	return redoc.Redoc(opt, handler)
}

func SetupRoutes(r *chi.Mux) {
	r.Get("/ping", s.StatusService)
	r.Route("/docs", func(r chi.Router) {
		r.Use(redocMiddleware)
		r.Get("/openapi.json", api.ServeOpenAPISpec)
	})
	r.Mount(PathPrefix(), apiRouter())
}

func apiRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/openapi.json", api.ServeOpenAPISpec)
	r.Group(func(r chi.Router) {
		r.Use(identity.EnforceIdentity)
		r.Use(middleware.AccountMiddleware)
		r.Route("/sources", func(r chi.Router) {
			r.Get("/", s.ListSources)
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", s.GetSource)
			})
		})

		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", s.ListAccounts)
			r.Route("/{ID}", func(r chi.Router) {
				r.Get("/", s.GetAccount)
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
			r.Post("/", s.CreateReservation)
		})

		r.Route("/instance_types", func(r chi.Router) {
			r.Route("/{source_id}", func(r chi.Router) {
				r.Get("/", s.ListInstanceTypes)
			})
		})
	})

	return r
}
