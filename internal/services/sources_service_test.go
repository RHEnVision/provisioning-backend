package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sources "github.com/RHEnVision/provisioning-backend/internal/clients/sources"
	"github.com/RHEnVision/provisioning-backend/internal/clients/sources/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func TestListSourcesHandler(t *testing.T) {
	t.SkipNow()
	ctx := identity.WithTenant(t, context.Background())
	ctx = stubs.WithSourcesClientV2(ctx)

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources", nil)
	assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListSources)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var sources []sources.Source

	if err := json.NewDecoder(rr.Body).Decode(&sources); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	assert.Equal(t, 2, len(sources), "expected two sources in response json")
}
