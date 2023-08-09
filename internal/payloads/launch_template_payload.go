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

type LaunchTemplateListResponse struct {
	Data []*LaunchTemplateResponse `json:"data" yaml:"data"`
}

func (s *LaunchTemplateResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *LaunchTemplateResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *LaunchTemplateListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListLaunchTemplateResponse(sl []*clients.LaunchTemplate) render.Renderer {
	list := make([]*LaunchTemplateResponse, len(sl))
	for i, instanceType := range sl {
		list[i] = &LaunchTemplateResponse{
			ID:   instanceType.ID,
			Name: instanceType.Name,
		}
	}
	return &LaunchTemplateListResponse{Data: list}
}
