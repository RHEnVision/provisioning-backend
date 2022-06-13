package payloads

import (
	"net/http"

	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/go-chi/render"
)

type SourceResponse struct {
	*sources.Source
}

func (s *SourceResponse) Bind(_ *http.Request) error {
	return nil
}

func (s *SourceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewListSourcesResponse(sl *[]sources.Source) []render.Renderer {
	sList := *sl
	list := make([]render.Renderer, 0, len(sList))
	for i := range sList {
		list = append(list, &SourceResponse{Source: &sList[i]})
	}
	return list
}

func NewShowSourcesResponse(s *sources.Source) render.Renderer {
	return &SourceResponse{Source: s}
}
