package stubs

import (
	"context"
	"net/http"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type sourcesCtxKeyType string

var sourcesCtxKey sourcesCtxKeyType = "sources-interface"

type SourcesIntegrationStub struct {
	store           *[]sources.Source
	authentications *[]sources.AuthenticationRead
}
type SourcesClientStub struct {
	sources []*clients.Source
	auths   map[string]*clients.Authentication
}

func init() {
	// We are currently using SourcesClientStub
	clients.GetSourcesClient = getSourcesClient
}

// SourcesClient
func WithSourcesClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, sourcesCtxKey, &SourcesClientStub{auths: make(map[string]*clients.Authentication)})
	return ctx
}

func AddSource(ctx context.Context, provider models.ProviderType) (*clients.Source, error) {
	stub, err := getSourcesClientStub(ctx)
	if err != nil {
		return nil, err
	}
	switch provider {
	case models.ProviderTypeAWS:
		return stub.addAuth(ctx, clients.NewAuthentication("arn:aws:iam::230214684733:role/Test", provider))
	case models.ProviderTypeAzure:
		return stub.addAuth(ctx, clients.NewAuthentication("4b9d213f-712f-4d17-a483-8a10bbe9df3a", provider))
	case models.ProviderTypeGCP:
		return stub.addAuth(ctx, clients.NewAuthentication("test@org.com", provider))
	case models.ProviderTypeUnknown, models.ProviderTypeNoop:
		// not implemented
		return nil, ErrNotImplemented
	}

	return nil, ErrNotImplemented
}

func AddAuth(ctx context.Context, authentication *clients.Authentication) (*clients.Source, error) {
	stub, err := getSourcesClientStub(ctx)
	if err != nil {
		return nil, err
	}
	return stub.addAuth(ctx, authentication)
}

func getSourcesClient(ctx context.Context) (clients.Sources, error) {
	return getSourcesClientStub(ctx)
}

func getSourcesClientStub(ctx context.Context) (si *SourcesClientStub, err error) {
	var ok bool
	if si, ok = ctx.Value(sourcesCtxKey).(*SourcesClientStub); !ok {
		err = ErrContextRead
	}
	return si, err
}

func (stub *SourcesClientStub) addAuth(ctx context.Context, authentication *clients.Authentication) (*clients.Source, error) {
	id := strconv.Itoa(len(stub.sources) + 2) // starts at 2 as 1 is reserved - TODO migrate users of the implicit id = 1
	source := &clients.Source{
		ID:   id,
		Name: "source-" + id,
	}
	stub.auths[id] = authentication
	stub.sources = append(stub.sources, source)
	return source, nil
}

// Implementation

func (*SourcesClientStub) Ready(ctx context.Context) error {
	return nil
}

func (stub *SourcesClientStub) GetAuthentication(ctx context.Context, sourceId string) (*clients.Authentication, error) {
	if sourceId == "1" {
		return clients.NewAuthentication("arn:aws:iam::230214684733:role/Test", models.ProviderTypeAWS), nil
	}

	auth, ok := stub.auths[sourceId]
	if !ok {
		return nil, ErrSourceAuthenticationNotFound
	}
	return auth, nil
}

func (mock *SourcesClientStub) ListAllProvisioningSources(ctx context.Context) ([]*clients.Source, int, error) {
	TestSourceData := []*clients.Source{
		{
			ID:           "1",
			Name:         "source1",
			SourceTypeID: "1",
			Uid:          "5eebe172-7baa-4280-823f-19e597d091e9",
			Provider:     models.ProviderTypeAWS,
		},
		{
			ID:           "2",
			Name:         "source2",
			SourceTypeID: "2",
			Uid:          "31b5338b-685d-4056-ba39-d00b4d7f19cc",
			Provider:     models.ProviderTypeAWS,
		},
	}
	return TestSourceData, 2, nil
}

func (mock *SourcesClientStub) ListProvisioningSourcesByProvider(ctx context.Context, provider models.ProviderType) ([]*clients.Source, int, error) {
	TestSourceData := []*clients.Source{
		{
			ID:           "1",
			Name:         "source1",
			SourceTypeID: "1",
			Uid:          "5eebe172-7baa-4280-823f-19e597d091e9",
			Provider:     models.ProviderTypeAWS,
		},
		{
			ID:           "2",
			Name:         "source2",
			SourceTypeID: "1",
			Uid:          "31b5338b-685d-4056-ba39-d00b4d7f19cc",
			Provider:     models.ProviderTypeAWS,
		},
	}
	return TestSourceData, 2, nil
}

// APIClient
func WithSourcesIntegration(parent context.Context, init_store *[]sources.Source) context.Context {
	ctx := context.WithValue(parent, sourcesCtxKey, &SourcesIntegrationStub{store: init_store})
	return ctx
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
