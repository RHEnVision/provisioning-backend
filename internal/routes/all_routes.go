package routes

import (
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Get("/ping", s.StatusService)
	r.Mount(pathPrefix(), apiRouter())
}

func pathPrefix() string {
	return fmt.Sprintf("/api/%s", config.Application.Name)
}

func apiRouter() http.Handler {
	r := chi.NewRouter()
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

	return r
}
