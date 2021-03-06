package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dao_stubs "github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
)

func buildPkStore() []*models.Pubkey {
	return []*models.Pubkey{&models.Pubkey{
		ID:        1,
		AccountID: 2,
		Name:      "firstkey",
		Body:      "sha-rsa body",
	}, &models.Pubkey{
		ID:        2,
		AccountID: 4,
		Name:      "secondkey",
		Body:      "sha-rsa body",
	}}
}

func TestListPubkeysHandler(t *testing.T) {
	ctx := identity.WithIdentity(t, context.Background())
	ctx = dao_stubs.WithPubkeyDao(ctx, buildPkStore())

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/pubkeys", nil)
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
	ctx := dao_stubs.WithPubkeyDao(context.Background(), nil)

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

	storecnt, err := dao_stubs.PubkeyStubCount(ctx)
	assert.Nil(t, err, fmt.Sprintf("Error reading stub count: %v", err))
	assert.Equal(t, 1, storecnt, "Pubkey has not been Created through DAO")
}
