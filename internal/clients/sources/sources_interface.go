package sources

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var GetSourcesClient func(ctx context.Context) (SourcesIntegration, error)
var GetSourcesClientV2 func(ctx context.Context) (ClientV2, error)

type SourcesIntegration interface {
	GetProvisioningTypeId(ctx context.Context, reqEditors ...RequestEditorFn) (string, error)
	ShowSourceWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*ShowSourceResponse, error)
	ListApplicationTypeSourcesWithResponse(ctx context.Context, appTypeId ID, params *ListApplicationTypeSourcesParams, reqEditors ...RequestEditorFn) (*ListApplicationTypeSourcesResponse, error)
	ListSourceAuthenticationsWithResponse(ctx context.Context, id ID, params *ListSourceAuthenticationsParams, reqEditors ...RequestEditorFn) (*ListSourceAuthenticationsResponse, error)
	ShowApplicationWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*ShowApplicationResponse, error)
}

type ClientV2 interface {
	// ListProvisioningSources returns all sources that have provisioning credentials assigned
	ListProvisioningSources(ctx context.Context) (*[]clients.Source, error)
	// GetArn returns ARN associated with provisioning app for given sourceId
	GetArn(ctx context.Context, sourceId ID) (string, error)
	// GetProvisioningTypeId might not need exposing
	GetProvisioningTypeId(ctx context.Context) (string, error)
}
