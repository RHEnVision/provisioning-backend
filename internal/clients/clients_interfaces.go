package clients

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/config"
)

var GetSourcesClient func(ctx context.Context) (SourcesIntegration, error)

// TODO move this initialization to sources package once we define clear interface
func getSourcesClient(ctx context.Context) (SourcesIntegration, error) {
	//nolint:wrapcheck
	return sources.NewClientWithResponses(config.Sources.URL)
}
func init() {
	GetSourcesClient = getSourcesClient
}

type SourcesIntegration interface {
	GetProvisioningTypeId(ctx context.Context, reqEditors ...sources.RequestEditorFn) (string, error)
	ShowSourceWithResponse(ctx context.Context, id sources.ID, reqEditors ...sources.RequestEditorFn) (*sources.ShowSourceResponse, error)
	ListApplicationTypeSourcesWithResponse(ctx context.Context, appTypeId sources.ID, params *sources.ListApplicationTypeSourcesParams, reqEditors ...sources.RequestEditorFn) (*sources.ListApplicationTypeSourcesResponse, error)
	ListSourceAuthenticationsWithResponse(ctx context.Context, id sources.ID, params *sources.ListSourceAuthenticationsParams, reqEditors ...sources.RequestEditorFn) (*sources.ListSourceAuthenticationsResponse, error)
	ShowApplicationWithResponse(ctx context.Context, id sources.ID, reqEditors ...sources.RequestEditorFn) (*sources.ShowApplicationResponse, error)
}
