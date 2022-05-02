package routes

import (
	m "github.com/RHEnVision/provisioning-backend/internal/middleware"
	s "github.com/RHEnVision/provisioning-backend/internal/services"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {
	r.Route("/ssh_keys", func(r chi.Router) {
		r.Get("/", s.ListSshKeys)
		r.Post("/", s.CreateSShKey)
		r.Route("/{ID}", func(r chi.Router) {
			r.Use(m.SshKeyCtx)
			r.Get("/", s.GetSshKey)
			r.Put("/", s.UpdateSshKey)
			r.Delete("/", s.DeleteSshKey)
			r.Route("/resources", func(r chi.Router) {
				r.Use(middleware.Timeout(time.Second * 1))
				r.Get("/", s.ListSshKeyResources)
				r.Post("/", s.CreateSshKeyResource)
				r.Route("/{RID}", func(r chi.Router) {
					r.Use(m.SshKeyResourceCtx)
					r.Delete("/", s.DeleteSshKeyResource)
				})
			})
		})
	})
}
