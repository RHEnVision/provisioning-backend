package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPubkeyValid(t *testing.T) {
	pk := examplePubkey
	err := Validate(context.Background(), &pk)
	assert.Nil(t, err)
}

func TestPubkeyMissingName(t *testing.T) {
	pk := Pubkey{Body: examplePubkey.Body}
	err := Validate(context.Background(), &pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Name' Error:Field validation for 'Name' failed on the 'required' tag")
}

func TestPubkeyMissingBody(t *testing.T) {
	pk := Pubkey{Name: examplePubkey.Name}
	err := Validate(context.Background(), &pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Body' Error:Field validation for 'Body' failed on the 'required' tag")
}

func TestPubkeyInvalidBody(t *testing.T) {
	pk := Pubkey{Name: examplePubkey.Name, Body: "XXX"}
	err := Validate(context.Background(), &pk)
	assert.EqualError(t, err, "Key: 'Pubkey.Body' Error:Field validation for 'Body' failed on the 'sshPubkey' tag")
}
