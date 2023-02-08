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
	"go.opentelemetry.io/otel"
)

const TraceName = "github.com/EnVision/provisioning/internal/clients/http/image_builder"

type ibClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetImageBuilderClient = newImageBuilderClient
}

func logger(ctx context.Context) zerolog.Logger {
	return ctxval.Logger(ctx).With().Str("client", "ib").Logger()
}

func newImageBuilderClient(ctx context.Context) (clients.ImageBuilder, error) {
	c, err := NewClientWithResponses(config.ImageBuilder.URL, func(c *Client) error {
		c.Client = http.NewPlatformClient(ctx, config.ImageBuilder.Proxy.URL)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ibClient{client: c}, nil
}

func (c *ibClient) Ready(ctx context.Context) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "Ready")
	defer span.End()

	logger := logger(ctx)
	resp, err := c.client.GetReadiness(ctx, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		logger.Error().Err(err).Msgf("Readiness request failed for image builder: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	err = http.HandleHTTPResponses(ctx, resp.StatusCode)
	if err != nil {
		return fmt.Errorf("ready call: %w", err)
	}
	return nil
}

func (c *ibClient) GetAWSAmi(ctx context.Context, composeID string) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting AMI of image %v", composeID)

	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	logger.Trace().Msgf("Verifying AWS type")
	if imageStatus.Type != UploadTypesAws {
		return "", fmt.Errorf("%w: expected image type AWS", http.UnknownImageTypeErr)
	}
	uploadStatus, err := imageStatus.Options.AsAWSUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not an AWS status", http.UploadStatusErr)
	}
	return uploadStatus.Ami, nil
}

func (c *ibClient) GetAzureImageName(ctx context.Context, composeID string) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting Azure ID of image %v", composeID)

	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	logger.Trace().Msgf("Verifying Azure type")
	if imageStatus.Type != UploadTypesAzure {
		return "", fmt.Errorf("%w: expected image type Azure, got %s", http.UnknownImageTypeErr, imageStatus.Type)
	}
	uploadStatus, err := imageStatus.Options.AsAzureUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not an Azure status", http.UploadStatusErr)
	}
	return uploadStatus.ImageName, nil
}

func (c *ibClient) GetGCPImageName(ctx context.Context, composeID string) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting Name of image %v", composeID)

	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	logger.Trace().Msg("Verifying GCP type")
	if imageStatus.Type != UploadTypesGcp {
		return "", fmt.Errorf("%w: expected image type GCP", http.UnknownImageTypeErr)
	}
	uploadStatus, err := imageStatus.Options.AsGCPUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not a GCP status", http.UploadStatusErr)
	}

	return fmt.Sprintf("projects/%s/global/images/%s", uploadStatus.ProjectId, uploadStatus.ImageName), nil
}

func (c *ibClient) fetchImageStatus(ctx context.Context, composeID string) (*UploadStatus, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "fetchImageStatus")
	defer span.End()
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v", composeID)

	composeResp, err := c.checkCompose(ctx, composeID)
	if err != nil {
		cloneResp, err := c.checkClone(ctx, composeID)
		if err != nil {
			return nil, fmt.Errorf("could not find image neither in compose nor in clones: %w", err)
		}
		return cloneResp, nil
	}
	return composeResp, nil
}

func (c *ibClient) checkCompose(ctx context.Context, composeID string) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v from composes", composeID)

	resp, err := c.client.GetComposeStatusWithResponse(ctx, composeID, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch image status from image builder")
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
		logger.Warn().Msg("Compose status is not ready")
		return nil, http.ImageStatusErr
	}

	return resp.JSON200.ImageStatus.UploadStatus, nil
}

func (c *ibClient) checkClone(ctx context.Context, composeID string) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v from clones", composeID)

	resp, err := c.client.GetCloneStatusWithResponse(ctx, composeID, headers.AddImageBuilderIdentityHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch image status from image builder")
		return nil, fmt.Errorf("cannot get compose status: %w", err)
	}

	err = http.HandleHTTPResponses(ctx, resp.StatusCode())
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			return nil, fmt.Errorf("fetch image status call: %w", http.CloneNotFoundErr)
		}
		return nil, fmt.Errorf("fetch image status call: %w", err)
	}

	if ImageStatusStatus(resp.JSON200.Status) != ImageStatusStatusSuccess {
		logger.Warn().Msg("Clone status is not ready")
		return nil, http.ImageStatusErr
	}

	return resp.JSON200, nil
}
