package services

import (
	"context"
	"net/http"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func AddIdentityHeader(ctx context.Context, req *http.Request) error {
	req.Header.Set("x-rh-identity", identity.GetIdentityHeader(ctx))
	return nil
}
