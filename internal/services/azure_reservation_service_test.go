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
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/queue/stub"
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAzureReservationHandler(t *testing.T) {
	var json_data []byte
	sharedCtx := stubs.WithAccountDaoOne(context.Background())
	sharedCtx = identity.WithTenant(t, sharedCtx)
	sharedCtx = Clientstubs.WithSourcesClient(sharedCtx)
	sharedCtx = Clientstubs.WithImageBuilderClient(sharedCtx)
	sharedCtx = stubs.WithPubkeyDao(sharedCtx)
	pk := factories.NewPubkeyRSA()
	err := stubs.AddPubkey(sharedCtx, pk)
	require.NoError(t, err, "failed to generate pubkey")
	source, err := Clientstubs.AddSource(sharedCtx, models.ProviderTypeAzure)
	require.NoError(t, err, "failed to generate Azure source")

	t.Run("successful reservation with compose ID", func(t *testing.T) {
		ctx := stubs.WithReservationDao(sharedCtx)
		ctx = stub.WithEnqueuer(ctx)

		var err error
		values := map[string]interface{}{
			"source_id":      source.ID,
			"image_id":       "92ea98f8-7697-472e-80b1-7454fa0e7fa7",
			"resource_group": "testGroup",
			"amount":         1,
			"instance_size":  "Basic_A0",
			"pubkey_id":      pk.ID,
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

		assert.Equal(t, 1, len(stub.EnqueuedJobs(ctx)), "Expected exactly one job to be planned")
		assert.IsType(t, jobs.LaunchInstanceAzureTaskArgs{}, stub.EnqueuedJobs(ctx)[0].Args, "Unexpected type of arguments for the planned job")
		jobArgs := stub.EnqueuedJobs(ctx)[0].Args.(jobs.LaunchInstanceAzureTaskArgs)
		assert.Equal(t, "testGroup", jobArgs.ResourceGroupName)
		assert.Equal(t, "/subscriptions/4b9d213f-712f-4d17-a483-8a10bbe9df3a/resourceGroups/redhat-deployed/providers/Microsoft.Compute/images/composer-api-92ea98f8-7697-472e-80b1-7454fa0e7fa7", jobArgs.AzureImageID, "Expected translated image to real name - one from IB client stub")
	})

	t.Run("successful reservation with azure image name translated to full azure ID", func(t *testing.T) {
		ctx := stubs.WithReservationDao(sharedCtx)
		ctx = stub.WithEnqueuer(ctx)

		var err error
		values := map[string]interface{}{
			"source_id":      source.ID,
			"image_id":       "composer-api-92ea98f8-7697-472e-80b1-7454fa0e7fa7",
			"resource_group": "testGroup",
			"amount":         1,
			"instance_size":  "Basic_A0",
			"pubkey_id":      pk.ID,
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

		assert.Equal(t, 1, len(stub.EnqueuedJobs(ctx)), "Expected exactly one job to be planned")
		assert.IsType(t, jobs.LaunchInstanceAzureTaskArgs{}, stub.EnqueuedJobs(ctx)[0].Args, "Unexpected type of arguments for the planned job")
		jobArgs := stub.EnqueuedJobs(ctx)[0].Args.(jobs.LaunchInstanceAzureTaskArgs)
		assert.Equal(t, "testGroup", jobArgs.ResourceGroupName)
		assert.Equal(t, "/subscriptions/4b9d213f-712f-4d17-a483-8a10bbe9df3a/resourceGroups/testGroup/providers/Microsoft.Compute/images/composer-api-92ea98f8-7697-472e-80b1-7454fa0e7fa7", jobArgs.AzureImageID, "Expected translated image to real name - one from IB client stub")
	})

	t.Run("failed reservation with invalid location", func(t *testing.T) {
		ctx := stubs.WithReservationDao(sharedCtx)
		ctx = stub.WithEnqueuer(ctx)

		var err error
		values := map[string]interface{}{
			"source_id":      source.ID,
			"location":       "blank",
			"image_id":       "92ea98f8-7697-472e-80b1-7454fa0e7fa7",
			"resource_group": "testGroup",
			"amount":         1,
			"instance_size":  "Basic_A0",
			"pubkey_id":      pk.ID,
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

		assert.Contains(t, rr.Body.String(), "Unsupported location")
		require.Equal(t, http.StatusBadRequest, rr.Code, "Handler returned wrong status code")
	})
}
