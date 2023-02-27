package services_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	Clientstubs "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGCPReservationHandler(t *testing.T) {
	var err error
	var json_data []byte
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = Clientstubs.WithSourcesClient(ctx)
	ctx = Clientstubs.WithImageBuilderClient(ctx)
	ctx = stubs.WithReservationDao(ctx)
	ctx = stubs.WithPubkeyDao(ctx)
	pk := factories.NewPubkeyRSA()
	err = stubs.AddPubkey(ctx, pk)
	require.NoError(t, err, "failed to generate pubkey")
	source, err := Clientstubs.AddSource(ctx, models.ProviderTypeGCP)
	require.NoError(t, err, "failed to generate GCP source")

	values := map[string]interface{}{
		"source_id":    source.Id,
		"image_id":     "80967e7f-efef-4eee-85b0-bd4cef4c455d",
		"amount":       1,
		"zone":         "us-central1-a",
		"machine_type": "n1-standard-1",
		"pubkey_id":    pk.ID,
	}
	if json_data, err = json.Marshal(values); err != nil {
		t.Fatalf("unable to marshal values to json: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/reservations/gcp", bytes.NewBuffer(json_data))
	require.NoError(t, err, "failed to create request")
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.CreateGCPReservation)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	stubCount := stubs.GCPReservationStubCount(ctx)
	assert.Equal(t, 1, stubCount, "Reservation has not been Created through DAO")
}
