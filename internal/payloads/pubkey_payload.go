package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/models"

	"github.com/go-chi/render"
)

type PubkeyPayload struct {
	*models.Pubkey
}
type (
	PubkeyRequest  PubkeyPayload
	PubkeyResponse PubkeyPayload
)

func (p *PubkeyRequest) Bind(_ *http.Request) error {
	// Fingerprint is read-only field
	p.Fingerprint = ""
	return nil
}

func (p *PubkeyResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewPubkeyResponse(pubkey *models.Pubkey) render.Renderer {
	return &PubkeyResponse{Pubkey: pubkey}
}

func NewPubkeyListResponse(pubkeys []*models.Pubkey) []render.Renderer {
	list := make([]render.Renderer, len(pubkeys))
	for i, pubkey := range pubkeys {
		list[i] = &PubkeyResponse{Pubkey: pubkey}
	}
	return list
}
