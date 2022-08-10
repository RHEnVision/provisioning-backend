package headers

import (
	"context"
	"net/http"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func AddIdentityHeader(ctx context.Context, req *http.Request) error {
	req.Header.Set("x-rh-identity", identity.GetIdentityHeader(ctx))
	return nil
}

func AddBasicAuth(ctx context.Context, req *http.Request) error {
	req.Header.Set("authorization", "Basic <REPLACE WITH STAGE BASE64(user:password)>")
	return nil
}
