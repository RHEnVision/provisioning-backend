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

func TestCreateAzureReservationHandler(t *testing.T) {
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
	source, err := Clientstubs.AddSource(ctx, models.ProviderTypeAzure)
	require.NoError(t, err, "failed to generate Azure source")

	values := map[string]interface{}{
		"source_id":     source.Id,
		"image_id":      "92ea98f8-7697-472e-80b1-7454fa0e7fa7",
		"amount":        1,
		"instance_size": "Basic_A0",
		"pubkey_id":     pk.ID,
	}
	if json_data, err = json.Marshal(values); err != nil {
		t.Fatalf("unable to marshal values to json: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/reservations/azure", bytes.NewBuffer(json_data))
	require.NoError(t, err, "failed to create request")
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(services.CreateAzureReservation)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	stubCount := stubs.AzureReservationStubCount(ctx)
	assert.Equal(t, 1, stubCount, "Reservation has not been Created through DAO")
}
