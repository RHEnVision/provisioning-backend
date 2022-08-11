package clients

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

var ClientProxyProductionUseErr = errors.New("client proxy cannot be used in production mode")

type ProxyClient struct {
	ctx    context.Context
	url    *url.URL
	client *http.Client
}

func NewProxyDoer(ctx context.Context, URL string) (*ProxyClient, error) {
	proxyURL, err := url.Parse(URL)
	if err != nil {
		return nil, fmt.Errorf("cannot create proxy doer: %w", err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := ProxyClient{
		ctx:    ctx,
		url:    proxyURL,
		client: &http.Client{Transport: transport},
	}
	return &client, nil
}

func (c *ProxyClient) Do(req *http.Request) (*http.Response, error) {
	ctxval.Logger(c.ctx).Trace().Msgf("Proxy request to %s via %s", req.URL, c.url.String())
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot proxy request: %w", err)
	}
	return resp, nil
}
