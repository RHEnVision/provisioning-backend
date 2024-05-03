package headers

import (
	"context"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/redhatinsights/platform-go-middlewares/v2/identity"
	"github.com/rs/zerolog"
)

func addIdentityHeader(ctx context.Context, req *http.Request, issuerUrl, clientID, clientSecret string) error {
	if clientID != "" && clientSecret != "" {
		zerolog.Ctx(ctx).Warn().Msgf("Using service account authentication: %s", clientID)
		token, err := getToken(ctx, issuerUrl, clientID, clientSecret)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Fetching access token failed")
			return err
		}
		req.Header.Add("Authorization", "Bearer "+token)
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
	issuerUrl := config.Sources.Issuer
	clientID := config.Sources.ClientID
	clientSecret := config.Sources.ClientSecret
	return addIdentityHeader(ctx, req, issuerUrl, clientID, clientSecret)
}

func AddImageBuilderIdentityHeader(ctx context.Context, req *http.Request) error {
	issuerUrl := config.ImageBuilder.Issuer
	clientID := config.ImageBuilder.ClientID
	clientSecret := config.ImageBuilder.ClientSecret
	return addIdentityHeader(ctx, req, issuerUrl, clientID, clientSecret)
}

func AddRbacIdentityHeader(ctx context.Context, req *http.Request) error {
	issuerUrl := config.RBAC.Issuer
	clientID := config.RBAC.ClientID
	clientSecret := config.RBAC.ClientSecret
	return addIdentityHeader(ctx, req, issuerUrl, clientID, clientSecret)
}
