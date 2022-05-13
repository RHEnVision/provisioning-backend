package routes

import (
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", s.ListAccounts)
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", s.GetAccount)
		})
	})
}
