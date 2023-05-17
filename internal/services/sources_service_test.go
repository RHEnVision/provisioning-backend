package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	clientStub "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func TestListSourcesHandler(t *testing.T) {
	t.Run("without provider", func(t *testing.T) {
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ListSources)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		var result []sources.Source

		err = json.NewDecoder(rr.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response body")

		assert.Equal(t, 2, len(result), "expected two result in response json")
	})

	t.Run("with provider", func(t *testing.T) {
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources?provider=aws", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ListSources)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		var result []payloads.SourceResponse

		err = json.NewDecoder(rr.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response body")

		assert.Equal(t, 2, len(result), "expected two result in response json")
	})

	t.Run("with invalid provider", func(t *testing.T) {
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources?provider=ibm", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ListSources)
		handler.ServeHTTP(rr, req)
		require.Error(t, clients.UnknownProviderErr, "provider is not supported")
		require.Equal(t, http.StatusBadRequest, rr.Code, "bad request")
	})
}

func TestGetAzureSourceDetails(t *testing.T) {
	t.Run("returns Azure details", func(t *testing.T) {
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)
		ctx = clientStub.WithAzureClient(ctx)

		sourceStub, err := clientStub.AddSource(ctx, models.ProviderTypeAzure)
		require.NoError(t, err, "failed to add stubbed source")

		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("ID", sourceStub.ID)
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("/api/provisioning/sources/%s/upload_info", sourceStub.ID), nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetSourceUploadInfo)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		var result payloads.SourceUploadInfoResponse

		err = json.NewDecoder(rr.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response body")

		assert.Equal(t, models.ProviderTypeAzure.String(), result.Provider, "Provider was expected to be Azure")
		assert.Equal(t, 3, len(result.AzureInfo.ResourceGroups), "expected three resource groups in response json")
	})
}
