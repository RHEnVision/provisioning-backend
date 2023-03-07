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
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAWSReservationHandler(t *testing.T) {
	var json_data []byte
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = Clientstubs.WithSourcesClient(ctx)
	ctx = Clientstubs.WithImageBuilderClient(ctx)
	ctx = stubs.WithReservationDao(ctx)
	ctx = stubs.WithPubkeyDao(ctx)
	pk := factories.NewPubkeyRSA()
	err := stubs.AddPubkey(ctx, pk)
	require.NoError(t, err, "failed to generate pubkey")

	t.Run("successful reservation", func(t *testing.T) {
		var err error
		values := map[string]interface{}{
			"source_id":     "1",
			"image_id":      "2bc640f6-927a-404a-9594-5b2da7e06608",
			"amount":        1,
			"instance_type": "t1.micro",
			"pubkey_id":     pk.ID,
		}
		if json_data, err = json.Marshal(values); err != nil {
			t.Fatalf("unable to marshal values to json: %v", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/reservations/aws", bytes.NewBuffer(json_data))
		require.NoError(t, err, "failed to create request")
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.CreateAWSReservation)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

		stubCount := stubs.AWSReservationStubCount(ctx)
		assert.Equal(t, 1, stubCount, "Reservation has not been created through DAO")
	})

	t.Run("failed reservation with invalid region", func(t *testing.T) {
		var err error
		values := map[string]interface{}{
			"source_id":     "1",
			"image_id":      "2bc640f6-927a-404a-9594-5b2da7e06608",
			"amount":        1,
			"instance_type": "t1.micro",
			"region":        "blank",
			"pubkey_id":     pk.ID,
		}
		if json_data, err = json.Marshal(values); err != nil {
			t.Fatalf("unable to marshal values to json: %v", err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/reservations/aws", bytes.NewBuffer(json_data))
		require.NoError(t, err, "failed to create request")
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.CreateAWSReservation)
		handler.ServeHTTP(rr, req)

		assert.Contains(t, rr.Body.String(), "Unsupported region")
		require.Equal(t, http.StatusBadRequest, rr.Code, "Handler returned wrong status code")
	})
}
