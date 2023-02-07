package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type LaunchTemplateResponse struct {
	*clients.LaunchTemplate
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
		list[i] = &LaunchTemplateResponse{instanceType}
	}
	return list
}
