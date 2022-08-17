package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"

	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func TestListPubkeysHandler(t *testing.T) {
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = stubs.WithPubkeyDao(ctx)
	err := stubs.AddPubkey(ctx, &models.Pubkey{
		Name: factories.GetSequenceName("pubkey"),
		Body: factories.GenerateRSAPubKey(t),
	})
	if err != nil {
		t.Fatalf("failed to add stubbed key: %v", err)
	}
	err = stubs.AddPubkey(ctx, &models.Pubkey{
		Name: factories.GetSequenceName("pubkey"),
		Body: factories.GenerateRSAPubKey(t),
	})
	if err != nil {
		t.Fatalf("failed to add stubbed key: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/pubkeys", nil)
	if err != nil {
		t.Fatalf("Error creating a new request: %v", err)
	}
	assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListPubkeys)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected: %d. Got: %d.", http.StatusOK, status)
	}

	var pubkeys []models.Pubkey

	if err := json.NewDecoder(rr.Body).Decode(&pubkeys); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	assert.Equal(t, 2, len(pubkeys), "expected two pubkeys in response json")
}

func TestCreatePubkeyHandler(t *testing.T) {
	var err error
	var json_data []byte
	ctx := stubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = stubs.WithPubkeyDao(ctx)

	values := map[string]interface{}{
		"account_id": 1,
		"name":       "very cool key",
		"body":       "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
	}

	if json_data, err = json.Marshal(values); err != nil {
		t.Fatal("unable to marshal values to json")
	}
	req, err := http.NewRequestWithContext(ctx, "POST", "/api/provisioning/pubkeys", bytes.NewBuffer(json_data))
	assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))
	req.Header.Add("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreatePubkey)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	storecnt, err := stubs.PubkeyStubCount(ctx)
	assert.Nil(t, err, fmt.Sprintf("Error reading stub count: %v", err))
	assert.Equal(t, 1, storecnt, "Pubkey has not been Created through DAO")
}
