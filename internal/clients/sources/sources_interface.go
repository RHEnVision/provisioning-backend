package sources

import (
	"context"
)

var GetSourcesClient func(ctx context.Context) (SourcesIntegration, error)

type SourcesIntegration interface {
	ShowSourceWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*ShowSourceResponse, error)
	ListApplicationTypeSourcesWithResponse(ctx context.Context, appId ID, params *ListApplicationTypeSourcesParams, reqEditors ...RequestEditorFn) (*ListApplicationTypeSourcesResponse, error)
	ListSourceAuthenticationsWithResponse(ctx context.Context, id ID, params *ListSourceAuthenticationsParams, reqEditors ...RequestEditorFn) (*ListSourceAuthenticationsResponse, error)
}
