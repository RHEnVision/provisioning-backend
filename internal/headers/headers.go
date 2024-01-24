package headers

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/rs/zerolog"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func addIdentityHeader(ctx context.Context, req *http.Request, username, password string) error {
	if username != "" && password != "" {
		zerolog.Ctx(ctx).Warn().Msgf("Username/password authentication: %s", username)
		req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	} else {
		logger := zerolog.Ctx(ctx)
		logger.Trace().Str("identity", identity.GetIdentityHeader(ctx)).Msg("HTTP client identity set")
		req.Header.Set("X-RH-Identity", identity.GetIdentityHeader(ctx))
	}
	return nil
}

func AddEdgeRequestIdHeader(ctx context.Context, req *http.Request) error {
	reqId := logging.EdgeRequestId(ctx)
	if reqId != "" {
		req.Header.Set("X-Rh-Edge-Request-Id", reqId)
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

func AddRbacIdentityHeader(ctx context.Context, req *http.Request) error {
	username := config.RBAC.Username
	password := config.RBAC.Password
	return addIdentityHeader(ctx, req, username, password)
}
