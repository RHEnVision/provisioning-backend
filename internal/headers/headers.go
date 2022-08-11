package headers

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/redhatinsights/platform-go-middlewares/identity"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func addIdentityHeader(ctx context.Context, req *http.Request, username, password string) error {
	if username != "" && password != "" {
		ctxval.Logger(ctx).Warn().Msgf("Username/password authentication: %s", username)
		req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	} else {
		req.Header.Set("X-RH-Identity", identity.GetIdentityHeader(ctx))
	}
	return nil
}

func AddSourcesIdentityHeader(ctx context.Context, req *http.Request) error {
	username := config.Sources.Username
	password := config.Sources.Password
	return addIdentityHeader(ctx, req, username, password)
}

func AddImageBuilderIdentityHeader(ctx context.Context, req *http.Request) error {
	username := config.ImageBuilder.Username
	password := config.ImageBuilder.Password
	return addIdentityHeader(ctx, req, username, password)
}
