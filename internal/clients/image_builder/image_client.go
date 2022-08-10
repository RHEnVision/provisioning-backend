package image_builder

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/RHEnVision/provisioning-backend/internal/parsing"
)

type ImageBuilderClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetImageBuilderClient = newImageBuilderClient
}

func newImageBuilderClient(ctx context.Context) (clients.ImageBuilder, error) {
	proxiedClient, err := clients.GetProxiedClient()
	c, err := NewClientWithResponses(config.ImageBuilder.URL, WithHTTPClient(proxiedClient))
	if err != nil {
		return nil, err
	}
	return &ImageBuilderClient{client: c}, nil
}

func (c *ImageBuilderClient) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	ctxval.Logger(ctx).Info().Msgf("Getting AMI of image %v", composeID)
	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}
	if imageStatus.Type != UploadTypesAws {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Image is not AWS type")
		return "", BadImageTypeErr
	}
	awsStatus := imageStatus.Options.(AWSUploadStatus)
	return awsStatus.Ami, nil
}

func (c *ImageBuilderClient) fetchImageStatus(ctx context.Context, composeID string) (*UploadStatus, error) {
	ctxval.Logger(ctx).Info().Msgf("Fetching image status %v", composeID)
	resp, err := c.client.GetComposeStatusWithResponse(ctx, composeID, headers.AddBasicAuth)
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Failed to fetch image status from image builder")
		return nil, fmt.Errorf("cannot get compose status: %w", err)
	}
	statusCode := resp.StatusCode()
	if parsing.IsHTTPNotFound(statusCode) {
		return nil, ComposeNotFoundErr
	}
	if !parsing.IsHTTPStatus2xx(statusCode) {
		ctxval.Logger(ctx).Warn().Msgf("Image builder replied with unexpected status while fetching image status: %v", statusCode)
		return nil, ClientErr
	}
	ctxval.Logger(ctx).Info().Msgf("Fetching image status was finished %+v\n", resp.JSON200.ImageStatus.UploadStatus)
	err = verifyImage(resp.JSON200)
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Image status in not ready")
		return nil, err
	}
	return resp.JSON200.ImageStatus.UploadStatus, nil

}

func verifyImage(compose *ComposeStatus) error {
	if compose.ImageStatus.Status != ImageStatusStatusSuccess {
		return ImageStatusErr
	}
	return nil
}
