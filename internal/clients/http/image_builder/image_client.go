package image_builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/headers"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ibClient struct {
	client *ClientWithResponses
}

func init() {
	clients.GetImageBuilderClient = newImageBuilderClient
}

func logger(ctx context.Context) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("client", "ib").Logger()
}

func newImageBuilderClient(ctx context.Context) (clients.ImageBuilder, error) {
	return NewImageBuilderClientWithUrl(ctx, config.ImageBuilder.URL)
}

func NewImageBuilderClientWithUrl(ctx context.Context, url string) (clients.ImageBuilder, error) {
	c, err := NewClientWithResponses(url, func(c *Client) error {
		c.Client = http.NewPlatformClient(ctx, config.ImageBuilder.Proxy.URL)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ibClient{client: c}, nil
}

func (c *ibClient) Ready(ctx context.Context) error {
	ctx, span := telemetry.StartSpan(ctx, "Ready")
	defer span.End()

	logger := logger(ctx)
	resp, err := c.client.GetReadinessWithResponse(ctx, headers.AddImageBuilderIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Error().Err(err).Msg("Readiness request failed for image builder")
		return err
	}

	if resp == nil {
		return fmt.Errorf("ready call: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() > 299 {
		return fmt.Errorf("ready call: %w: %d", clients.ErrUnexpectedBackendResponse, resp.StatusCode())
	}

	return nil
}

func (c *ibClient) GetAWSAmi(ctx context.Context, composeUUID uuid.UUID, instanceType clients.InstanceType) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting AMI of compose ID %s", composeUUID.String())

	imageStatus, err := c.fetchImageStatus(ctx, composeUUID, instanceType)
	if err != nil {
		return "", err
	}
	if imageStatus == nil {
		return "", fmt.Errorf("%w: no image status", http.ErrImageStatus)
	}

	if imageStatus.Type != UploadTypesAws {
		return "", fmt.Errorf("%w: expected image type AWS", http.ErrUnknownImageType)
	}
	uploadStatus, err := imageStatus.Options.AsAWSUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not an AWS status", http.ErrUploadStatus)
	}

	logger.Info().Msgf("Translated compose ID %s to AMI %s", composeUUID, uploadStatus.Ami)

	return uploadStatus.Ami, nil
}

func (c *ibClient) GetAzureImageInfo(ctx context.Context, composeUUID uuid.UUID, instanceType clients.InstanceType) (string, string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting Azure ID of image %v", composeUUID.String())

	composeStatus, err := c.getComposeStatus(ctx, composeUUID)
	if err != nil {
		return "", "", err
	}
	if composeStatus == nil {
		logger.Warn().Msg("Compose status is not ready")
		return "", "", fmt.Errorf("getting azure id: %w", http.ErrImageStatus)
	}

	logger.Trace().Msgf("Verifying Azure type")
	if composeStatus.ImageStatus.UploadStatus == nil {
		return "", "", fmt.Errorf("%w: upload status is nil", http.ErrUploadStatus)
	}
	if composeStatus.ImageStatus.UploadStatus.Type != UploadTypesAzure {
		return "", "", fmt.Errorf("%w: expected image type Azure, got %s", http.ErrUnknownImageType, composeStatus.ImageStatus.UploadStatus.Type)
	}
	if len(composeStatus.Request.ImageRequests) < 1 {
		logger.Error().Msg(http.ErrImageRequestNotFound.Error())
		return "", "", http.ErrImageRequestNotFound
	}

	imageArch, archErr := clients.MapArchitectures(ctx, string(composeStatus.Request.ImageRequests[0].Architecture))
	if archErr != nil || imageArch != instanceType.Architecture {
		return "", "", http.ErrImageArchInvalid
	}

	uploadOptions, err := composeStatus.ImageStatus.UploadStatus.Options.AsAzureUploadStatus()
	if err != nil {
		return "", "", fmt.Errorf("%w: not an Azure status", http.ErrUploadStatus)
	}

	azureUploadRequest, err := composeStatus.Request.ImageRequests[0].UploadRequest.Options.AsAzureUploadRequestOptions()
	if err != nil {
		return "", "", fmt.Errorf("failed to decode Azure upload request from IB: %w", err)
	}
	return azureUploadRequest.ResourceGroup, uploadOptions.ImageName, nil
}

func (c *ibClient) GetGCPImageName(ctx context.Context, composeUUID uuid.UUID, instanceType clients.InstanceType) (string, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Getting Google image id of compose %s", composeUUID)

	imageStatus, err := c.fetchImageStatus(ctx, composeUUID, instanceType)
	if err != nil {
		return "", err
	}

	if imageStatus == nil {
		return "", fmt.Errorf("%w: no image status", http.ErrImageStatus)
	}

	if imageStatus.Type != UploadTypesGcp {
		return "", fmt.Errorf("%w: expected image type GCP", http.ErrUnknownImageType)
	}
	uploadStatus, err := imageStatus.Options.AsGCPUploadStatus()
	if err != nil {
		return "", fmt.Errorf("%w: not a GCP status", http.ErrUploadStatus)
	}

	result := fmt.Sprintf("projects/%s/global/images/%s", uploadStatus.ProjectId, uploadStatus.ImageName)
	logger.Info().Msgf("Translated compose ID %s to image name %s", composeUUID, result)

	return result, nil
}

func (c *ibClient) fetchImageStatus(ctx context.Context, composeUUID uuid.UUID, instanceType clients.InstanceType) (*UploadStatus, error) {
	ctx, span := telemetry.StartSpan(ctx, "fetchImageStatus")
	defer span.End()
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v", composeUUID.String())

	uploadStatus, err := c.checkCompose(ctx, composeUUID, instanceType)
	if err != nil {
		if errors.Is(err, clients.ErrUnexpectedBackendResponse) {
			uploadStatus, err = c.checkClone(ctx, composeUUID)
			if err != nil {
				return nil, fmt.Errorf("could not find image neither in compose nor in clones: %w", err)
			}
		} else {
			return nil, fmt.Errorf("image compose is not launchable: %w", err)
		}
	}
	return uploadStatus, nil
}

func (c *ibClient) getComposeStatus(ctx context.Context, composeUUID uuid.UUID) (*ComposeStatus, error) {
	logger := logger(ctx)

	resp, err := c.client.GetComposeStatusWithResponse(ctx, composeUUID, headers.AddImageBuilderIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch image status from image builder")
		return nil, fmt.Errorf("cannot get compose status: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("cannot get compose status: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("fetch image status call: %w", clients.ErrUnexpectedBackendResponse)
	}

	return resp.JSON200, nil
}

// checkCompose validates whether the image identified by composeUUID is already successfully built.
// It also checks whether the target instance type architecture matches the image architecture.
func (c *ibClient) checkCompose(ctx context.Context, composeUUID uuid.UUID, instanceType clients.InstanceType) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v from composes", composeUUID)

	composeStatus, err := c.getComposeStatus(ctx, composeUUID)
	if err != nil {
		return nil, err
	}

	if composeStatus == nil || composeStatus.ImageStatus.Status != ImageStatusStatusSuccess {
		logger.Warn().Msg("Compose status is not ready")
		return nil, http.ErrImageStatus
	}

	if len(composeStatus.Request.ImageRequests) == 1 {
		imageArch, archErr := clients.MapArchitectures(ctx, string(composeStatus.Request.ImageRequests[0].Architecture))
		if archErr != nil || imageArch != instanceType.Architecture {
			return nil, http.ErrImageArchInvalid
		}
	}

	return composeStatus.ImageStatus.UploadStatus, nil
}

func (c *ibClient) checkClone(ctx context.Context, composeUUID uuid.UUID) (*UploadStatus, error) {
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching image status %v from clones", composeUUID)

	resp, err := c.client.GetCloneStatusWithResponse(ctx, composeUUID, headers.AddImageBuilderIdentityHeader, headers.AddEdgeRequestIdHeader)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to fetch image status from image builder")
		return nil, fmt.Errorf("cannot get compose status: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("cannot get compose status: empty response: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("fetch image status call: %w", clients.ErrUnexpectedBackendResponse)
	}

	if resp.JSON200.Status != UploadStatusStatusSuccess {
		logger.Warn().Msgf("Clone status (%s) is not ready", resp.JSON200.Status)
		return nil, http.ErrImageStatus
	}

	return resp.JSON200, nil
}
