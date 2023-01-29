package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type InstanceDescriptionResponse struct {
	*clients.InstanceDescription
}

func (s *InstanceDescriptionResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *InstanceDescriptionResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListInstanceDescriptionResponse(instances []*clients.InstanceDescription) []render.Renderer {
	list := make([]render.Renderer, len(instances))
	for i, instance := range instances {
		list[i] = &InstanceDescriptionResponse{instance}
	}
	return list
}
