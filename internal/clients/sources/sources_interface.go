package sources

import (
	"context"
)

var GetSourcesClient func(ctx context.Context) (SourcesIntegration, error)

type SourcesIntegration interface {
	ShowSourceWithResponse(ctx context.Context, id ID, reqEditors ...RequestEditorFn) (*ShowSourceResponse, error)
	ListSourcesWithResponse(ctx context.Context, params *ListSourcesParams, reqEditors ...RequestEditorFn) (*ListSourcesResponse, error)
}
