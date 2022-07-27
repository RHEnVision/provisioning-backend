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
	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
)

func buildSourcesStore() *[]sources.Source {
	var TestSourceData = []sources.Source{
		{
			Id:           ptr.String("1"),
			Name:         ptr.String("source1"),
			SourceTypeId: ptr.String("1"),
			Uid:          ptr.String("5eebe172-7baa-4280-823f-19e597d091e9"),
		},
		{
			Id:           ptr.String("2"),
			Name:         ptr.String("source2"),
			SourceTypeId: ptr.String("2"),
			Uid:          ptr.String("31b5338b-685d-4056-ba39-d00b4d7f19cc"),
		},
	}
	return &TestSourceData
}

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
