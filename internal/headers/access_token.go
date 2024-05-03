package headers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/rs/zerolog"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
}

const (
	scopes     = "openid api.iam.service_accounts"
	grant_type = "client_credentials"
)

// https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/token" -d "scope=openid api.iam.service_accounts
func getToken(ctx context.Context, issuerUrl, clientId, clientSecret string) (string, error) {
	provider, err := oidc.NewProvider(ctx, issuerUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch oidc provider info: %w", err)
	}

	data := url.Values{}
	data.Add("grant_type", grant_type)
	data.Add("scope", scopes)
	data.Add("client_id", clientId)
	data.Add("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, provider.Endpoint().TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to form request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request a token: %w", err)
	}
	defer res.Body.Close()

	token := tokenResponse{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse a token response: %w", err)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		return "", fmt.Errorf("failed to parse a token response: %w", err)
	}

	zerolog.Ctx(ctx).Debug().Msgf("Fetched access token: %s", token.AccessToken)

	return token.AccessToken, nil
}
