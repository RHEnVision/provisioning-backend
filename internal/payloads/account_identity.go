package payloads

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/go-chi/render"
)

type AccountIdentityResponse struct {
	*clients.AccountIdentity
}

func NewAccountIdentityResponse(awsAccountId string) render.Renderer {
	return &AccountIdentityResponse{
		&clients.AccountIdentity{
			AWSDetails: &clients.AccountDetailsAWS{
				AccountID: awsAccountId,
			},
		},
	}
}

func (s *AccountIdentityResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
