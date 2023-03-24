package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type AccountIdentityResponse struct {
	*clients.AccountIdentity
}

func NewAccountIdentityResponse(accDetails *clients.AccountDetailsAWS) render.Renderer {
	return &AccountIdentityResponse{
		&clients.AccountIdentity{
			AWSDetails: accDetails,
		},
	}
}

func (s *AccountIdentityResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
