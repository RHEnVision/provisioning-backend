package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	ibClientStub "github.com/RHEnVision/provisioning-backend/internal/clients/image_builder/stubs"
	sourcesClientStub "github.com/RHEnVision/provisioning-backend/internal/clients/sources/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/assert"
)

func TestCreateAWSReservationHandler(t *testing.T) {
	var err error
	var json_data []byte
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = sourcesClientStub.WithSourcesClient(ctx)
	ctx = ibClientStub.WithImageBuilderClient(ctx)
	ctx = stubs.WithReservationDao(ctx)
	ctx = stubs.WithPubkeyDao(ctx)
	pk := models.Pubkey{Name: "new", Body: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC8w6DONv1qn3IdgxSpkYOClq7oe7davWFqKVHPbLoS6+dFInru7gdEO5byhTih6+PwRhHv/b1I+Mtt5MDZ8Sv7XFYpX/3P/u5zQiy1PkMSFSz0brRRUfEQxhXLW97FJa7l+bej2HJDt7f9Gvcj+d/fNWC9Z58/GX11kWk4SIXaKotkN+kWn54xGGS7Zvtm86fP59Srt6wlklSsG8mZBF7jVUjyhAgm/V5gDFb2/6jfiwSb2HyJ9/NbhLkWNdwrvpdGZqQlYhnwTfEZdpwizW/Mj3MxP5O31HN45aE0wog0UeWY4gvTl4Ogb6kescizAM6pCff3RBslbFxLdOO7cR17 lzap+rsakey@redhat.com"}
	err = stubs.AddPubkey(ctx, &pk)
	assert.Nil(t, err, fmt.Sprintf("Error GeneratePubkey: %v", err))

	values := map[string]interface{}{
		"source_id":     1,
		"image_id":      "2bc640f6-927a-404a-9594-5b2da7e06608",
		"amount":        1,
		"instance_type": "t1.micro",
		"pubkey_id":     pk.ID,
	}
	if json_data, err = json.Marshal(values); err != nil {
		t.Fatal("unable to marshal values to json")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/reservations/aws", bytes.NewBuffer(json_data))
	req.Header.Add("Content-Type", "application/json")
	assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateAWSReservation)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	storecnt, err := stubs.ReservationStubCount(ctx)
	assert.Nil(t, err, fmt.Sprintf("Error reading stub count: %v", err))
	assert.Equal(t, 1, storecnt, "Reservation has not been Created through DAO")
}
