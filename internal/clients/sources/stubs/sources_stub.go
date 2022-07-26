package stubs

import (
	"context"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/aws/smithy-go/ptr"
)

type sourcesCtxKeyType string

var sourcesCtxKey sourcesCtxKeyType = "sources-interface"

type SourcesIntegrationStub struct {
	sources         *[]sources.Source
	authentications *[]sources.AuthenticationRead
}

func init() {
	sources.GetSourcesClient = getSourcesClientStub
}

type contextReadError struct{}

func (m *contextReadError) Error() string {
	return "failed to find or convert dao stored in testing context"
}

func WithSourcesIntegration(parent context.Context, sources *[]sources.Source, authentications *[]sources.AuthenticationRead) context.Context {
	ctx := context.WithValue(parent, sourcesCtxKey, &SourcesIntegrationStub{sources: sources, authentications: authentications})
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
	lst := *mock.sources
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
			Data: mock.sources,
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

func (mock *SourcesIntegrationStub) FetchARN(ctx context.Context, sourceId string) (string, error) {
	return "arn:aws:iam::230934684733:role/Test", nil
}

func (mock *SourcesIntegrationStub) FilterSourceAuthentications(authentications *[]sources.AuthenticationRead) (sources.AuthenticationRead, error) {
	return sources.AuthenticationRead{
		ResourceType: (*sources.AuthenticationReadResourceType)(ptr.String("Application")),
		Name:         ptr.String("test"),
		ResourceId:   ptr.String("1"),
	}, nil
}
