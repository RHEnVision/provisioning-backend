package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dao_stubs "github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/test/identity"
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
	ctx := context.Background()
	ctx = identity.WithIdentity(t, ctx)
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
