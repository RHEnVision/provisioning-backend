package payloads

import (
	"net/http"

	"github.com/go-chi/render"
)

type PermissionsResponse struct {
	Valid           bool     `json:"valid"`
	MissingEntities []string `json:"missing_entities,omitempty"`
}

func (s *PermissionsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewPermissionsResponse(sl []string) render.Renderer {
	response := PermissionsResponse{
		Valid:           len(sl) == 0,
		MissingEntities: sl,
	}
	return &response
}
