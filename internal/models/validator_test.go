package models_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/stretchr/testify/assert"
)

func TestPubkeyValid(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	err := models.Validate(context.Background(), pk)
	assert.Nil(t, err)
}

func TestPubkeyMissingName(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	pk.Name = ""
	err := models.Validate(context.Background(), pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Name' Error:Field validation for 'Name' failed on the 'required' tag")
}

func TestPubkeyMissingBody(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	pk.Body = ""
	err := models.Validate(context.Background(), pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Body' Error:Field validation for 'Body' failed on the 'required' tag")
}

func TestPubkeyInvalidBody(t *testing.T) {
	pk := factories.NewPubkeyED25519()
	pk.Body = "xxx"
	err := models.Validate(context.Background(), pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Body' Error:Field validation for 'Body' failed on the 'sshPubkey' tag")
}
