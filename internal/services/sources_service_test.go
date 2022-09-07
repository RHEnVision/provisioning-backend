package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/require"

	"github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"
	clientStub "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func TestListSourcesHandler(t *testing.T) {
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
}
