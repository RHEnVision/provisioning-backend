package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var examplePubkey = Pubkey{
	Name: "lzap-ed25519-2021",
	Body: "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
}

func TestPubkeyFingerprintGeneration(t *testing.T) {
	pk := examplePubkey
	err := Transform(context.Background(), &pk)
	if assert.NoError(t, err) {
		assert.Equal(t, "SHA256:gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk", pk.Fingerprint, "fingerprint must be set")
	}
}
