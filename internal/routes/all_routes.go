package routes

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/api"
	"github.com/RHEnVision/provisioning-backend/internal/routes/helpers"
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/go-chi/chi/v5"
	redoc "github.com/go-openapi/runtime/middleware"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func redocMiddleware(handler http.Handler) http.Handler {
	opt := redoc.RedocOpts{
		SpecURL: fmt.Sprintf("%s/openapi.json", helpers.PathPrefix()),
	}
	return redoc.Redoc(opt, handler)
}

func SetupRoutes(r *chi.Mux) {
	r.Get("/ping", s.StatusService)
	r.Route("/docs", func(r chi.Router) {
		r.Use(redocMiddleware)
		r.Get("/openapi.json", api.ServeOpenAPISpec)
	})
	r.Mount(helpers.PathPrefix(), apiRouter())
}

func apiRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/openapi.json", api.ServeOpenAPISpec)
	r.Group(func(r chi.Router) {
		r.Use(identity.EnforceIdentity)
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
				// this is temporary until we implement tasks
				r.Get("/upload_aws", s.UploadPubkeyResourceAWS)
				r.Get("/delete_all", s.DeleteAllPubkeyResources)
			})
		})
	})

	return r
}
