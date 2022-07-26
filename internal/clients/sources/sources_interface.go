package sources

import (
	"context"
)

var GetSourcesClient func(ctx context.Context) (SourcesIntegration, error)

type SourcesIntegration interface {
	GetProvisioningTypeId(ctx context.Context, reqEditors ...RequestEditorFn) (string, error)
	ShowSourceWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*ShowSourceResponse, error)
	ListApplicationTypeSourcesWithResponse(ctx context.Context, appTypeId ID, params *ListApplicationTypeSourcesParams, reqEditors ...RequestEditorFn) (*ListApplicationTypeSourcesResponse, error)
	ListSourceAuthenticationsWithResponse(ctx context.Context, id ID, params *ListSourceAuthenticationsParams, reqEditors ...RequestEditorFn) (*ListSourceAuthenticationsResponse, error)
	ShowApplicationWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*ShowApplicationResponse, error)
	FetchARN(ctx context.Context, sourceId string) (string, error)
	FilterSourceAuthentications(authentications *[]AuthenticationRead) (AuthenticationRead, error)
}
