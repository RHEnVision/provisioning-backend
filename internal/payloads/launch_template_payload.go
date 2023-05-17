package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

// See clients.LaunchTemplate
type LaunchTemplateResponse struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func (s *LaunchTemplateResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *LaunchTemplateResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListLaunchTemplateResponse(sl []*clients.LaunchTemplate) []render.Renderer {
	list := make([]render.Renderer, len(sl))
	for i, instanceType := range sl {
		list[i] = &LaunchTemplateResponse{
			ID:   instanceType.ID,
			Name: instanceType.Name,
		}
	}
	return list
}
