package middleware

// This is copied from github.com/redhatinsights/platform-go-middlewares/identity
// to allow for logging into clowdwatch, logging should be fixed in the global library tho.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/redhatinsights/platform-go-middlewares/identity"
	"github.com/rs/zerolog"
)

var (
	InvalidOrMissingOrgErr = errors.New("X-Rh-Identity header has an invalid or missing org_id")
	MissingTypeErr         = errors.New("X-Rh-Identity header is missing type")
)

func doError(ctx context.Context, w http.ResponseWriter, code int, reason string) {
	logger := zerolog.Ctx(ctx)
	logger.Warn().Msgf("Failed to enforce the Identity header: %s", reason)
	http.Error(w, reason, code)
}

// EnforceIdentity extracts the X-Rh-Identity header and places the contents into the
// request context.  If the Identity is invalid, the request will be aborted.
func EnforceIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())
		rawHeaders := r.Header["X-Rh-Identity"]

		// must have an x-rh-id header
		if len(rawHeaders) != 1 {
			doError(r.Context(), w, 400, "missing X-Rh-Identity header")
			return
		}

		// must be able to base64 decode header
		idRaw, err := base64.StdEncoding.DecodeString(rawHeaders[0])
		if err != nil {
			doError(r.Context(), w, 400, "unable to b64 decode x-rh-identity header")
			return
		}

		var jsonData identity.XRHID
		err = json.Unmarshal(idRaw, &jsonData)
		if err != nil {
			logger.Warn().Err(err).Msg("unable to unmarshal X-Rh-Identity header")
			doError(r.Context(), w, 400, "X-Rh-Identity header does not contain valid JSON")
			return
		}

		topLevelOrgIDFallback(&jsonData)

		logger.Debug().RawJSON("user", []byte(idRaw)).Msg("Enforcing identity")
		err = checkHeader(r.Context(), &jsonData, w)
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), identity.Key, jsonData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkHeader(ctx context.Context, id *identity.XRHID, w http.ResponseWriter) error {
	if (id.Identity.Type == "Associate" || id.Identity.Type == "X509") && id.Identity.AccountNumber == "" {
		return nil
	}

	if id.Identity.OrgID == "" && id.Identity.Internal.OrgID == "" {
		doError(ctx, w, 400, "X-Rh-Identity header has an invalid or missing org_id")
		return InvalidOrMissingOrgErr
	}

	if id.Identity.Type == "" {
		doError(ctx, w, 400, "X-Rh-Identity header is missing type")
		return MissingTypeErr
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
