package services_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	identity2 "github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/payloads"
	"github.com/RHEnVision/provisioning-backend/internal/services"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetReservationDetail(t *testing.T) {
	t.Run("Generic reservation", func(t *testing.T) {
		var err error
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = identity.WithTenant(t, ctx)
		ctx = stubs.WithPubkeyDao(ctx)
		ctx = stubs.WithReservationDao(ctx)
		pk := &models.Pubkey{
			Name: factories.SeqNameWithPrefix("pubkey"),
			Body: factories.GenerateRSAPubKey(t),
		}
		err = stubs.AddPubkey(ctx, pk)
		require.NoError(t, err, "failed to add stubbed key")

		detail := &models.AWSDetail{
			Region:       "us-east-1",
			InstanceType: "t1.micro",
			Amount:       2,
			PowerOff:     true,
		}
		reservation := &models.AWSReservation{
			PubkeyID: pk.ID,
			SourceID: "1",
			ImageID:  "ami-random",
			Detail:   detail,
		}
		reservation.AccountID = identity2.AccountId(ctx)
		reservation.Status = "Created"
		reservation.Provider = models.ProviderTypeAWS
		reservation.Steps = 2

		err = stubs.AddAWSReservation(ctx, reservation)
		require.NoError(t, err, "failed to create stub reservation")

		rctx := chi.NewRouteContext()
		ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
		rctx.URLParams.Add("ID", "1")
		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/v1/reservations/1", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.GetReservationDetail)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Wrong status code")

		var response payloads.GenericReservationResponsePayload
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err, "failed to decode response body")

		assert.Equal(t, int(models.ProviderTypeAWS), response.Provider, "expected provider to be AWS in parsed json")
	})
}
