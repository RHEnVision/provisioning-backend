package payloads

import (
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"net/http"

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

func NewPubkeyResponse(account *models.Pubkey) render.Renderer {
	return &PubkeyResponse{Pubkey: account}
}

func NewPubkeyListResponse(accounts []*models.Pubkey) []render.Renderer {
	list := make([]render.Renderer, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &PubkeyResponse{Pubkey: a})
	}
	return list
}
