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

	r.Route("/pubkeys", func(r chi.Router) {
		r.Post("/", s.CreatePubkey)
		r.Get("/", s.ListPubkeys)
		r.Route("/{ID}", func(r chi.Router) {
			r.Get("/", s.GetPubkey)
			r.Delete("/", s.DeletePubkey)
		})
	})
}
