package payloads

import (
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"net/http"

	"github.com/go-chi/render"
)

type SSHKeyPayload struct {
	*models.SSHKey
}
type SSHKeyRequest SSHKeyPayload
type SSHKeyResponse SSHKeyPayload

func (p *SSHKeyRequest) Bind(r *http.Request) error {
	return nil
}

func (p *SSHKeyResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewSshKeyResponse(sshKey *models.SSHKey) render.Renderer {
	return &SSHKeyResponse{SSHKey: sshKey}
}

func NewSSHKeyListResponse(sshKeys []*models.SSHKey) []render.Renderer {
	var list []render.Renderer
	for _, k := range sshKeys {
		list = append(list, &SSHKeyResponse{SSHKey: k})
	}
	return list
}
