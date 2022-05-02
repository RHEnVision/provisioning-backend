package middleware

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	m "github.com/RHEnVision/provisioning-backend/internal/models"
	p "github.com/RHEnVision/provisioning-backend/internal/payloads"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func parseInt64(r *http.Request, param string) (int64, error) {
	i, err := strconv.Atoi(chi.URLParam(r, param))
	if err != nil {
		return 0, err
	} else {
		return int64(i), nil
	}
}

func SshKeyCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int64
		var sshKey *m.SSHKey
		var err error

		if id, err = parseInt64(r, "ID"); err == nil {
			sshKey, err = m.FindSSHKey(r.Context(), db.DB, id)
			if err != nil {
				render.Render(w, r, p.ErrNotFound)
				return
			}
		} else if err != nil {
			render.Render(w, r, p.ErrParamParsingError)
			return
		} else {
			render.Render(w, r, p.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), ctxval.SshKeyCtxKey, sshKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SshKeyFromCtx(ctx context.Context) *m.SSHKey {
	return ctx.Value(ctxval.SshKeyCtxKey).(*m.SSHKey)
}

func SshKeyResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int64
		var sshKey *m.SSHKeyResource
		var err error

		if id, err = parseInt64(r, "RID"); err == nil {
			sshKey, err = m.FindSSHKeyResource(r.Context(), db.DB, id)
			if err != nil {
				render.Render(w, r, p.ErrNotFound)
				return
			}
		} else if err != nil {
			render.Render(w, r, p.ErrParamParsingError)
			return
		} else {
			render.Render(w, r, p.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), ctxval.SshKeyResourceCtxKey, sshKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SshKeyResourceFromCtx(ctx context.Context) *m.SSHKeyResource {
	return ctx.Value(ctxval.SshKeyResourceCtxKey).(*m.SSHKeyResource)
}
