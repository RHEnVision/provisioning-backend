package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type SourceID struct {
	SourceId string `json:"source_id"`
}
type SourceResponse struct {
	*clients.Source
}

func (s *SourceResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *SourceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListSourcesResponse(sourceList []*clients.Source) []render.Renderer {
	list := make([]render.Renderer, len(sourceList))
	for i, source := range sourceList {
		list[i] = &SourceResponse{Source: source}
	}
	return list
}
