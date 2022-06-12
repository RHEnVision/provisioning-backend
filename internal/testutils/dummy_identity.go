package testutils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// accountNumber to be used in the tests.
const accountNumber = "12345"

// orgId to be used in the tests.
const orgId = "abc-org-id"

var xRhId = identity.XRHID{
	Identity: identity.Identity{
		AccountNumber: accountNumber,
		OrgID:         orgId,
	},
}

func AddIdentityHeader(t *testing.T, req *http.Request) *http.Request {
	req.Header.Add("X-Rh-Identity", setUpValidIdentity(t))
	return req
}

func WithIdentity(t *testing.T, ctx context.Context) context.Context {
	return context.WithValue(ctx, identity.Key, xRhId)
}

// setUpValidIdentity returns a base64 encoded valid identity.
func setUpValidIdentity(t *testing.T) string {
	jsonIdentity, err := json.Marshal(xRhId)
	if err != nil {
		t.Errorf(`could not marshal test identity to JSON: %s`, err)
	}

	base64Identity := base64.StdEncoding.EncodeToString(jsonIdentity)

	return base64Identity
}
