package payloads

import (
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"net/http"

	"github.com/go-chi/render"
)

type SSHKeyResourcePayload struct {
	*models.SSHKeyResource
}
type SSHKeyResourceRequest SSHKeyResourcePayload
type SSHKeyResourceResponse SSHKeyResourcePayload

func (p *SSHKeyResourceRequest) Bind(r *http.Request) error {
	return nil
}

func (p *SSHKeyResourceResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewSSHKeyResourceListResponse(sshKeys []*models.SSHKeyResource) []render.Renderer {
	var list []render.Renderer
	for _, k := range sshKeys {
		list = append(list, &SSHKeyResourceResponse{SSHKeyResource: k})
	}
	return list
}
