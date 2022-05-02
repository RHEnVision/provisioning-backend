package services

import (
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	m "github.com/RHEnVision/provisioning-backend/internal/models"
	p "github.com/RHEnVision/provisioning-backend/internal/payloads"
	"net/http"

	"github.com/go-chi/render"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func CreateSShKey(w http.ResponseWriter, r *http.Request) {
	data := &p.SSHKeyRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, p.ErrInvalidRequest(err))
		return
	}

	sshKey := data.SSHKey
	sshKey.InsertP(r.Context(), db.DB, boil.Infer())

	render.Status(r, http.StatusCreated)
	render.Render(w, r, p.NewSshKeyResponse(sshKey))
}

func ListSshKeys(w http.ResponseWriter, r *http.Request) {
	logger := ContextLogger(r)
	logger.Info().Msg("Listing ssh keys")
	keys := m.SSHKeys().AllP(r.Context(), db.DB)
	if err := render.RenderList(w, r, p.NewSSHKeyListResponse(keys)); err != nil {
		render.Render(w, r, p.ErrRender(err))
		return
	}
}

func GetSshKey(w http.ResponseWriter, r *http.Request) {
	sshKey := r.Context().Value(ctxval.SshKeyCtxKey).(*m.SSHKey)
	sshKey.SSHKeyResources()
	render.Render(w, r, p.NewSshKeyResponse(sshKey))
}

func DeleteSshKey(w http.ResponseWriter, r *http.Request) {
	sshKey := r.Context().Value(ctxval.SshKeyCtxKey).(*m.SSHKey)
	rows := sshKey.DeleteP(r.Context(), db.DB)
	if rows == 1 {
		render.Render(w, r, p.NewSshKeyResponse(sshKey))
	} else {
		render.Render(w, r, p.ErrDeleteError)
	}
}

func UpdateSshKey(w http.ResponseWriter, r *http.Request) {
	existing := r.Context().Value(ctxval.SshKeyCtxKey).(*m.SSHKey)
	data := &p.SSHKeyRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, p.ErrInvalidRequest(err))
		return
	}

	updated := data.SSHKey
	existing.Body = updated.Body

	rows := existing.UpdateP(r.Context(), db.DB, boil.Infer())
	if rows == 1 {
		render.Render(w, r, p.NewSshKeyResponse(existing))
	} else {
		render.Render(w, r, p.ErrDeleteError)
	}
}
