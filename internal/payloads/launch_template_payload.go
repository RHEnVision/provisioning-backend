package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/go-chi/render"
)

// See clients.LaunchTemplate
type LaunchTemplateResponse struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type LaunchTemplateListResponse struct {
	Data     []*LaunchTemplateResponse `json:"data" yaml:"data"`
	Metadata page.Metadata             `json:"metadata" yaml:"metadata"`
}

func (s *LaunchTemplateListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListLaunchTemplateResponse(sl []*clients.LaunchTemplate, meta *page.Metadata) render.Renderer {
	list := make([]*LaunchTemplateResponse, len(sl))
	for i, tmpl := range sl {
		list[i] = &LaunchTemplateResponse{
			ID:   tmpl.ID,
			Name: tmpl.Name,
		}
	}
	return &LaunchTemplateListResponse{Data: list, Metadata: *meta}
}
