// Copyright Red Hat

package common

import (
	"fmt"
	"net/http"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// Get account from http request header
func GetAccount(r *http.Request) (string, error) {
	if r.Context().Value(identity.Key) != nil {
		ident := identity.Get(r.Context())
		if ident.Identity.AccountNumber != "" {
			return ident.Identity.AccountNumber, nil
		}
	}
	return "", fmt.Errorf("cannot find account number")

}