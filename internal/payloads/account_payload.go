package payloads

import (
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"net/http"

	"github.com/go-chi/render"
)

type AccountPayload struct {
	*models.Account
}
type AccountRequest AccountPayload
type AccountResponse AccountPayload

func (p *AccountRequest) Bind(_ *http.Request) error {
	return nil
}

func (p *AccountResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewAccountResponse(account *models.Account) render.Renderer {
	return &AccountResponse{Account: account}
}

func NewAccountListResponse(accounts []*models.Account) []render.Renderer {
	list := make([]render.Renderer, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &AccountResponse{Account: a})
	}
	return list
}
