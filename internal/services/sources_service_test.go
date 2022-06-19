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

var (
	uid1 = "5eebe172-7baa-4280-823f-19e597d091e9"
	uid2 = "31b5338b-685d-4056-ba39-d00b4d7f19cc"
)

func buildSourcesStore() *[]sources.Source {
	source1 := "source1"
	id1 := "1"
	type1 := "1"
	source2 := "source2"
	id2 := "2"
	type2 := "2"

	var TestSourceData = []sources.Source{
		{
			Id:           &id1,
			Name:         &source1,
			SourceTypeId: &type1,
			Uid:          &uid1,
		},
		{
			Id:           &id2,
			Name:         &source2,
			SourceTypeId: &type2,
			Uid:          &uid2,
		},
	}
	return &TestSourceData
}
func buildSource() *[]sources.Source {
	source1 := "source1"
	id1 := "1"
	type1 := "1"
	var TestSourceData = []sources.Source{
		{
			Id:           &id1,
			Name:         &source1,
			SourceTypeId: &type1,
			Uid:          &uid1,
		},
	}
	return &TestSourceData
}
func TestListSourcesHandler(t *testing.T) {
	ctx := identity.WithIdentity(t, context.Background())
	ctx = stubs.WithSourcesIntegration(ctx, buildSourcesStore())

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

func TestShowSourceHandler(t *testing.T) {
	ctx := identity.WithIdentity(t, context.Background())
	ctx = stubs.WithSourcesIntegration(ctx, buildSource())

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/sources/1", nil)
	assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetSource)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var s sources.Source

	if err := json.NewDecoder(rr.Body).Decode(&s); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	assert.Equal(t, "1", *s.Id, "expected source with id = 1")
}
