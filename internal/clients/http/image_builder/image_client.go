package image_builder

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
)

type ImageBuilderClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetImageBuilderClient = newImageBuilderClient
}

func newImageBuilderClient(ctx context.Context) (clients.ImageBuilder, error) {
	c, err := NewClientWithResponses(config.ImageBuilder.URL, func(c *Client) error {
		if config.ImageBuilder.Proxy.URL != "" {
			var client HttpRequestDoer
			if config.Features.Environment != "development" {
				return clients.ClientProxyProductionUseErr
			}
			client, err := clients.NewProxyDoer(ctx, config.ImageBuilder.Proxy.URL)
			if err != nil {
				return fmt.Errorf("cannot create proxy doer: %w", err)
			}
			if config.RestEndpoints.TraceData {
				client = clients.NewLoggingDoer(ctx, client)
			}
			c.Client = client
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ImageBuilderClient{client: c}, nil
}

func (c *ImageBuilderClient) Ready(ctx context.Context) error {
	resp, err := c.client.GetReadiness(ctx, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		ctxval.Logger(ctx).Error().Err(err).Msgf("Readiness request failed for image builder: %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	if !http.IsHTTPStatus2xx(resp.StatusCode) {
		ctxval.Logger(ctx).Warn().Msgf("Readiness response from image builder: %d", resp.StatusCode)
		return ClientErr
	}
	return nil
}

func (c *ImageBuilderClient) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	ctxval.Logger(ctx).Info().Msgf("Getting AMI of image %v", composeID)
	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}
	ctxval.Logger(ctx).Debug().Msgf("Verifying AWS type")
	if imageStatus.Type != UploadTypesAws {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Image is not AWS type")
		return "", BadImageTypeErr
	}
	awsStatus, ok := imageStatus.Options.(map[string]interface{})
	if !ok {
		return "", BadImageTypeErr
	}
	return awsStatus["ami"].(string), nil
}

func (c *ImageBuilderClient) fetchImageStatus(ctx context.Context, composeID string) (*UploadStatus, error) {
	ctxval.Logger(ctx).Info().Msgf("Fetching image status %v", composeID)
	resp, err := c.client.GetComposeStatusWithResponse(ctx, composeID, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Failed to fetch image status from image builder")
		return nil, fmt.Errorf("cannot get compose status: %w", err)
	}
	statusCode := resp.StatusCode()
	if http.IsHTTPNotFound(statusCode) {
		return nil, ComposeNotFoundErr
	}
	if !http.IsHTTPStatus2xx(statusCode) {
		ctxval.Logger(ctx).Warn().Msgf("Image builder replied with unexpected status while fetching image status: %v", statusCode)
		return nil, ClientErr
	}

	err = verifyImage(ctx, resp.JSON200)
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Image status in not ready")
		return nil, err
	}
	return resp.JSON200.ImageStatus.UploadStatus, nil

}

func verifyImage(ctx context.Context, compose *ComposeStatus) error {
	ctxval.Logger(ctx).Debug().Msgf("Verifying image is ready to use")
	if compose.ImageStatus.Status != ImageStatusStatusSuccess {
		return ImageStatusErr
	}
	return nil
}
