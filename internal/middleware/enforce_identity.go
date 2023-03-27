package middleware

// This is copied from github.com/redhatinsights/platform-go-middlewares/identity
// to allow for logging into clowdwatch, logging should be fixed in the global library tho.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/go-chi/render"
	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
)

var ErrIdentity = errors.New("identity error")

// EnforceIdentity extracts the X-Rh-Identity header and places the contents into the
// request context. If the Identity is invalid, the request will be aborted.
func EnforceIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())
		rawHeaders := r.Header["X-Rh-Identity"]

		// must have an x-rh-id header
		if len(rawHeaders) != 1 {
			errRender := render.Render(w, r, payloads.NewMissingIdentityError(r.Context(), "missing X-Rh-Identity header", ErrIdentity))
			if errRender != nil {
				logger.Warn().Err(errRender).Msg("Cannot render permission middleware error")
			}
			return
		}

		// must be able to base64 decode header
		idRaw, err := base64.StdEncoding.DecodeString(rawHeaders[0])
		if err != nil {
			errRender := render.Render(w, r, payloads.NewMissingIdentityError(r.Context(), "unable to b64 decode X-Rh-Identity header", ErrIdentity))
			if errRender != nil {
				logger.Warn().Err(errRender).Msg("Cannot render permission middleware error")
			}
			return
		}

		var jsonData identity.XRHID
		err = json.Unmarshal(idRaw, &jsonData)
		if err != nil {
			errRender := render.Render(w, r, payloads.NewMissingIdentityError(r.Context(), "X-Rh-Identity header does not contain valid JSON", ErrIdentity))
			if errRender != nil {
				logger.Warn().Err(errRender).Msg("Cannot render permission middleware error")
			}
			return
		}

		topLevelOrgIDFallback(&jsonData)

		logger.Debug().RawJSON("user", idRaw).Msg("Enforcing identity")
		err = checkHeader(r.Context(), &jsonData, w, r)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), identity.Key, jsonData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkHeader(ctx context.Context, id *identity.XRHID, w http.ResponseWriter, r *http.Request) error {
	if (id.Identity.Type == "Associate" || id.Identity.Type == "X509") && id.Identity.AccountNumber == "" {
		return nil
	}

	if id.Identity.OrgID == "" && id.Identity.Internal.OrgID == "" {
		errRender := render.Render(w, r, payloads.NewMissingIdentityError(r.Context(), "X-Rh-Identity header has an invalid or missing org_id", ErrIdentity))
		if errRender != nil {
			zerolog.Ctx(ctx).Warn().Err(errRender).Msg("Cannot render permission middleware error")
		}
		return ErrIdentity
	}

	if id.Identity.Type == "" {
		errRender := render.Render(w, r, payloads.NewMissingIdentityError(r.Context(), "X-Rh-Identity header is missing type", ErrIdentity))
		if errRender != nil {
			zerolog.Ctx(ctx).Warn().Err(errRender).Msg("Cannot render permission middleware error")
		}
		return ErrIdentity
	}

	return nil
}

// if org_id is not defined at the top level, use the internal one
// https://issues.redhat.com/browse/RHCLOUD-17717
func topLevelOrgIDFallback(identity *identity.XRHID) {
	if identity.Identity.OrgID == "" && identity.Identity.Internal.OrgID != "" {
		identity.Identity.OrgID = identity.Identity.Internal.OrgID
	}
}
