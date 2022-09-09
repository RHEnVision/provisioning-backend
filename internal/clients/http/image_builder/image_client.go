package image_builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/rs/zerolog"
)

// TODO This should have been not exported
type ImageBuilderClient struct {
	client *ClientWithResponses
	logger zerolog.Logger
}

func init() {
	clients.GetImageBuilderClient = newImageBuilderClient
}

func newImageBuilderClient(ctx context.Context) (clients.ImageBuilder, error) {
	logger := ctxval.Logger(ctx).With().Str("client", "ib").Logger()
	c, err := NewClientWithResponses(config.ImageBuilder.URL, func(c *Client) error {
		if config.ImageBuilder.Proxy.URL != "" {
			var client HttpRequestDoer
			if config.Features.Environment != "development" {
				return http.ClientProxyProductionUseErr
			}
			client, err := http.NewProxyDoer(ctx, config.ImageBuilder.Proxy.URL)
			if err != nil {
				return fmt.Errorf("cannot create proxy doer: %w", err)
			}
			if config.RestEndpoints.TraceData {
				client = http.NewLoggingDoer(ctx, client)
			}
			c.Client = client
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ImageBuilderClient{client: c, logger: logger}, nil
}

func (c *ImageBuilderClient) Ready(ctx context.Context) error {
	resp, err := c.client.GetReadiness(ctx, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		c.logger.Error().Err(err).Msgf("Readiness request failed for image builder: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		return fmt.Errorf("ready call: %w", err)
	}
	return nil
}

func (c *ImageBuilderClient) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	c.logger.Trace().Msgf("Getting AMI of image %v", composeID)

	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	c.logger.Trace().Msgf("Verifying AWS type")
	if imageStatus.Type != UploadTypesAws {
		return "", fmt.Errorf("%w: expected image type AWS", http.UnknownImageTypeErr)
	}
	ami, ok := imageStatus.Options.(map[string]interface{})["ami"]
	if !ok {
		return "", http.AmiNotFoundInStatusErr
	}
	return ami.(string), nil
}

func (c *ImageBuilderClient) fetchImageStatus(ctx context.Context, composeID string) (*UploadStatus, error) {
	c.logger.Trace().Msgf("Fetching image status %v", composeID)

	resp, err := c.client.GetComposeStatusWithResponse(ctx, composeID, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		c.logger.Warn().Err(err).Msg("Failed to fetch image status from image builder")
		return nil, fmt.Errorf("cannot get compose status: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			return nil, fmt.Errorf("fetch image status call: %w", http.ComposeNotFoundErr)
		}
		return nil, fmt.Errorf("fetch image status call: %w", err)
	}

	if resp.JSON200.ImageStatus.Status != ImageStatusStatusSuccess {
		c.logger.Warn().Msg("Image status in not ready")
		return nil, http.ImageStatusErr
	}
	return resp.JSON200.ImageStatus.UploadStatus, nil
}
