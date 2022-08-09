package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/models"

	"github.com/go-chi/render"
)

type AccountPayload struct {
	// Required auto-generated PK.
	ID int64 `json:"id"`

	// Organization ID. Required.
	OrgID string `json:"org_id"`

	// EBS account number. Can be NULL but not blank.
	AccountNumber *string `json:"account_number"`
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
	return &AccountResponse{ID: account.ID, OrgID: account.OrgID, AccountNumber: SqlNullToStringPtr(account.AccountNumber)}
}

func NewAccountListResponse(accounts []*models.Account) []render.Renderer {
	list := make([]render.Renderer, 0, len(accounts))
	for _, a := range accounts {
		list = append(list, &AccountResponse{ID: a.ID, OrgID: a.OrgID, AccountNumber: SqlNullToStringPtr(a.AccountNumber)})
	}
	return list
}
