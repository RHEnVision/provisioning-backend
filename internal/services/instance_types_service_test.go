package services_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	clientStub "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func TestListInstanceTypes(t *testing.T) {

	t.Run("with region", func(t *testing.T) {
		var names []string
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = clientStub.WithSourcesClient(ctx)
		ctx = clientStub.WithEC2Client(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources/1/instance_types?region=us-east-1", nil)
		assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.ListInstanceTypes)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
		}

		var result []clients.InstanceType

		if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}

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
		assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.ListInstanceTypes)
		handler.ServeHTTP(rr, req)

		assert.Error(t, payloads.ParamMissingError{})
	})

}
