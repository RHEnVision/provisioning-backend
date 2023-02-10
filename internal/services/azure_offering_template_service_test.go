package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/stretchr/testify/require"
)

func TestAzureOfferingTemplate(t *testing.T) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/v1/azure_offering_template", nil)
	require.NoError(t, err, "failed to create request")
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.AzureOfferingTemplate)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	var data map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &data)
	require.NoError(t, err, "failed to parse the template from response")

	// Check template has expected key - just simple check for now
	_, ok := data["parameters"]
	require.True(t, ok, "the rendered template does not seem correct")
}
