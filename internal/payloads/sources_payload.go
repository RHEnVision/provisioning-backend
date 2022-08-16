package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type SourceID struct {
	SourceId int64 `json:"source_id"`
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

func NewListSourcesResponse(sl *[]clients.Source) []render.Renderer {
	sList := *sl
	list := make([]render.Renderer, 0, len(sList))
	for i := range sList {
		list = append(list, &SourceResponse{Source: &sList[i]})
	}
	return list
}
