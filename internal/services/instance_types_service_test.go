package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	clientStub "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInstanceTypesHandler(t *testing.T) {

	t.Run("with region", func(t *testing.T) {
		var names []string
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)
		ctx = clientStub.WithEC2Client(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources/1/instance_types?region=us-east-1", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.ListInstanceTypes)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		var result []clients.InstanceType

		err = json.NewDecoder(rr.Body).Decode(&result)
		require.NoError(t, err, "failed to decode response body")

		assert.Equal(t, 3, len(result), "expected three result in response json")
		for _, it := range result {
			names = append(names, it.Name.String())
		}
		assert.Contains(t, names, "a1.2xlarge", "expected result to contain a1.2xlarge instance type")
		assert.Contains(t, names, "c5.xlarge", "expected result to contain c5.xlarge instance type")
	})

	t.Run("without region", func(t *testing.T) {
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)
		ctx = clientStub.WithEC2Client(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources/1/instance_types", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.ListInstanceTypes)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusBadRequest, rr.Code, "Handler returned wrong status code")

		assert.Contains(t, rr.Body.String(), "missing parameter")
	})

}
