package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/rbac"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
)

var ErrMissingExtraPermission = errors.New("missing permission")

// removeEmptyStrings returns a copy of a slice without empty strings
func removeEmptyStrings(s []string) []string {
	r := make([]string, 0, len(s))
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// CheckPermissionAndRender can be used to perform an extra permission check that is more detailed than the one
// performed by the middleware. Do not use this function as the only permission check, permissions should always
// be enforced via middleware as a bare minimum.
func CheckPermissionAndRender(w http.ResponseWriter, r *http.Request, permission string, resources ...string) error {
	resource := strings.Join(removeEmptyStrings(resources), ".")
	if !rbac.Acl(r.Context()).IsAllowed(resource, permission) {
		permErr := fmt.Errorf("%w: %s on %s", ErrMissingExtraPermission, permission, resource)
		renderError(w, r, payloads.NewMissingPermissionError(r.Context(), resource, permission, permErr))
		return permErr
	}

	return nil
}
