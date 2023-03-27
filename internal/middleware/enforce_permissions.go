package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/rbac"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

var (
	ErrEnforceIdentityFirst = errors.New("missing identity in the context, use EnforceIdentity before EnforcePermissions")
	ErrMissingPermission    = errors.New("missing permission")
)

// EnforcePermissions enforces permissions via RBAC service. It requires that identity is present
// in the context, make sure to chain EnforceIdentity middleware before this one.
func EnforcePermissions(resource, permission string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(r.Context())
			rbacClient := clients.GetRbacClient(r.Context())

			if identity.Identity(r.Context()).Identity.OrgID == "" {
				panic(ErrEnforceIdentityFirst)
			}

			acl, err := rbacClient.GetPrincipalAccess(r.Context())
			logger.Trace().Str("acl_resource", resource).Str("acl_permission", permission).
				Msgf("Checking permission '%s' on '%s'", permission, resource)

			if err != nil {
				aclErr := fmt.Errorf("unable to get ACL: %w", err)
				errRender := render.Render(w, r, payloads.NewClientError(r.Context(), aclErr))
				if errRender != nil {
					logger.Warn().Err(errRender).Msg("Cannot render permission middleware error")
				}
				return
			}

			if !acl.IsAllowed(resource, permission) {
				permErr := fmt.Errorf("%w: %s on %s", ErrMissingPermission, permission, resource)
				errRender := render.Render(w, r, payloads.NewMissingPermissionError(r.Context(), resource, permission, permErr))
				if errRender != nil {
					logger.Warn().Err(errRender).Msg("Cannot render permission middleware error")
				}
				return
			}

			ctx := rbac.WithAcl(r.Context(), acl)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
