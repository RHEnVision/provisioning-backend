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
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
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
		var doer HttpRequestDoer
		doer, err := telemetry.HTTPClient(ctx, config.StringToURL(ctx, config.ImageBuilder.Proxy.URL))
		if err != nil {
			return fmt.Errorf("cannot HTTP client: %w", err)
		}
		if config.RestEndpoints.TraceData {
			doer = http.NewLoggingDoer(ctx, doer)
		}
		c.Client = doer
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
	ami, ok := imageStatus.Options.(map[string]interface{})["ami"]
	if !ok {
		return "", http.AmiNotFoundInStatusErr
	}
	return ami.(string), nil
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
	imageName, ok := imageStatus.Options.(map[string]interface{})["image_name"]
	if !ok {
		return "", fmt.Errorf("%w: image name was not found", http.NameNotFoundInStatusErr)
	}

	projectID, ok := imageStatus.Options.(map[string]interface{})["project_id"]
	if !ok {
		return "", fmt.Errorf("%w: project id was not found", http.IDNotFoundInStatusErr)
	}

	return fmt.Sprintf("projects/%s/global/images/%s", projectID, imageName.(string)), nil
}

func (c *ibClient) fetchImageStatus(ctx context.Context, composeID string) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v", composeID)

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
		logger.Warn().Msg("Image status in not ready")
		return nil, http.ImageStatusErr
	}
	return resp.JSON200.ImageStatus.UploadStatus, nil
}
