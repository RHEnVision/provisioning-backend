package stubs

import (
	"context"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients/sources"
)

type sourcesCtxKeyType string

var sourcesCtxKey sourcesCtxKeyType = "sources-interface"

type SourcesIntegrationStub struct {
	store           *[]sources.Source
	authentications *[]sources.AuthenticationRead
}

func init() {
	sources.GetSourcesClient = getSourcesClientStub
}

type contextReadError struct{}

func (m *contextReadError) Error() string {
	return "failed to find or convert dao stored in testing context"
}

func WithSourcesIntegration(parent context.Context, init_store *[]sources.Source) context.Context {
	ctx := context.WithValue(parent, sourcesCtxKey, &SourcesIntegrationStub{store: init_store})
	return ctx
}

func getSourcesClientStub(ctx context.Context) (si sources.SourcesIntegration, err error) {
	var ok bool
	if si, ok = ctx.Value(sourcesCtxKey).(*SourcesIntegrationStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *SourcesIntegrationStub) GetProvisioningTypeId(ctx context.Context, reqEditors ...sources.RequestEditorFn) (string, error) {
	return "11", nil
}

func (mock *SourcesIntegrationStub) ShowSourceWithResponse(ctx context.Context, id sources.ID, reqEditors ...sources.RequestEditorFn) (*sources.ShowSourceResponse, error) {
	lst := *mock.store
	return &sources.ShowSourceResponse{
		JSON200: &lst[0],
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}
func (mock *SourcesIntegrationStub) ListApplicationTypeSourcesWithResponse(ctx context.Context, appTypeId sources.ID, params *sources.ListApplicationTypeSourcesParams, reqEditors ...sources.RequestEditorFn) (*sources.ListApplicationTypeSourcesResponse, error) {
	return &sources.ListApplicationTypeSourcesResponse{
		JSON200: &sources.SourcesCollection{
			Data: mock.store,
		},
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}

func (mock *SourcesIntegrationStub) ListSourceAuthenticationsWithResponse(ctx context.Context, sourceId sources.ID, params *sources.ListSourceAuthenticationsParams, reqEditors ...sources.RequestEditorFn) (*sources.ListSourceAuthenticationsResponse, error) {
	return &sources.ListSourceAuthenticationsResponse{
		JSON200: &sources.AuthenticationsCollection{
			Data: mock.authentications,
		},
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}

func (mock *SourcesIntegrationStub) ShowApplicationWithResponse(ctx context.Context, appId sources.ID, reqEditors ...sources.RequestEditorFn) (*sources.ShowApplicationResponse, error) {
	return &sources.ShowApplicationResponse{
		JSON200: &sources.Application{},
		HTTPResponse: &http.Response{
			StatusCode: 200,
		},
	}, nil
}
