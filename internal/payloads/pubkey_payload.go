package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/models"

	"github.com/go-chi/render"
)

type PubkeyPayload struct {
	*models.Pubkey
}
type PubkeyRequest PubkeyPayload
type PubkeyResponse PubkeyPayload

func (p *PubkeyRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *PubkeyResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewPubkeyResponse(pubkey *models.Pubkey) render.Renderer {
	return &PubkeyResponse{Pubkey: pubkey}
}

func NewPubkeyListResponse(pubkeys []*models.Pubkey) []render.Renderer {
	list := make([]render.Renderer, 0, len(pubkeys))
	for _, pubkey := range pubkeys {
		list = append(list, &PubkeyResponse{Pubkey: pubkey})
	}
	return list
}
