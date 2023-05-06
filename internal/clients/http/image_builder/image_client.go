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
	"github.com/google/uuid"
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
	resp, err := c.client.GetReadiness(ctx, headers.AddImageBuilderIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Error().Err(err).Msg("Readiness request failed for image builder")
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
	logger.Trace().Str("compose_id", composeID).Msgf("Getting AMI of compose ID %v", composeID)

	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	if imageStatus.Type != UploadTypesAws {
		return "", fmt.Errorf("%w: expected image type AWS", http.UnknownImageTypeErr)
	}
	uploadStatus, err := imageStatus.Options.AsAWSUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not an AWS status", http.UploadStatusErr)
	}

	logger.Info().Str("compose_id", composeID).Str("ami", uploadStatus.Ami).
		Msgf("Translated compose ID %s to AMI %s", composeID, uploadStatus.Ami)

	return uploadStatus.Ami, nil
}

func (c *ibClient) GetAzureImageID(ctx context.Context, composeID string) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting Azure ID of image %v", composeID)

	composeStatus, err := c.getComposeStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	logger.Trace().Msgf("Verifying Azure type")
	if composeStatus.ImageStatus.UploadStatus.Type != UploadTypesAzure {
		return "", fmt.Errorf("%w: expected image type Azure, got %s", http.UnknownImageTypeErr, composeStatus.ImageStatus.UploadStatus.Type)
	}
	if len(composeStatus.Request.ImageRequests) < 1 {
		logger.Error().Msg(http.ImageRequestNotFoundErr.Error())
		return "", http.ImageRequestNotFoundErr
	}

	uploadOptions, err := composeStatus.ImageStatus.UploadStatus.Options.AsAzureUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not an Azure status", http.UploadStatusErr)
	}

	azureUploadRequest, err := composeStatus.Request.ImageRequests[0].UploadRequest.Options.AsAzureUploadRequestOptions()
	if err != nil {
		return "", fmt.Errorf("failed to decode Azure upload request from IB: %w", err)
	}
	return fmt.Sprintf("/resourceGroups/%s/providers/Microsoft.Compute/images/%s", azureUploadRequest.ResourceGroup, uploadOptions.ImageName), nil
}

func (c *ibClient) GetGCPImageName(ctx context.Context, composeID string) (string, error) {
	logger := logger(ctx)
	logger.Trace().Str("compose_id", composeID).Msgf("Getting Google image id of compose %s", composeID)

	imageStatus, err := c.fetchImageStatus(ctx, composeID)
	if err != nil {
		return "", err
	}

	if imageStatus.Type != UploadTypesGcp {
		return "", fmt.Errorf("%w: expected image type GCP", http.UnknownImageTypeErr)
	}
	uploadStatus, err := imageStatus.Options.AsGCPUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not a GCP status", http.UploadStatusErr)
	}

	result := fmt.Sprintf("projects/%s/global/images/%s", uploadStatus.ProjectId, uploadStatus.ImageName)
	logger.Info().Str("compose_id", composeID).Str("ami", result).
		Msgf("Translated compose ID %s to AMI %s", composeID, result)

	return result, nil
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

func (c *ibClient) getComposeStatus(ctx context.Context, composeID string) (*ComposeStatus, error) {
	logger := logger(ctx)

	composeUUID, err := uuid.Parse(composeID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse UUID: %w", err)
	}

	resp, err := c.client.GetComposeStatusWithResponse(ctx, composeUUID, headers.AddImageBuilderIdentityHeader, headers.AddEdgeRequestIdHeader)
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

	return resp.JSON200, nil
}

func (c *ibClient) checkCompose(ctx context.Context, composeID string) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v from composes", composeID)

	composeStatus, err := c.getComposeStatus(ctx, composeID)
	if err != nil {
		return nil, err
	}

	if composeStatus.ImageStatus.Status != ImageStatusStatusSuccess {
		logger.Warn().Msg("Compose status is not ready")
		return nil, http.ImageStatusErr
	}

	return composeStatus.ImageStatus.UploadStatus, nil
}

func (c *ibClient) checkClone(ctx context.Context, composeID string) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v from clones", composeID)

	composeUUID, err := uuid.Parse(composeID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse UUID: %w", err)
	}

	resp, err := c.client.GetCloneStatusWithResponse(ctx, composeUUID, headers.AddImageBuilderIdentityHeader, headers.AddEdgeRequestIdHeader)
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
