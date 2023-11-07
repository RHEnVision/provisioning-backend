package models_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFingerprintGeneration(t *testing.T) {
	type test struct {
		name        string
		pubkey      *models.Pubkey
		fingerprint string
	}

	tests := []test{
		{"ed25519", factories.NewPubkeyED25519(), "gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk="},
		{"rsa", factories.NewPubkeyRSA(), "ENShRe/0uDLSw9c+7tc9PxkD/p4blyB/DTgBSIyTAJY="},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			err := models.Transform(context.Background(), td.pubkey)
			require.NoError(t, err)
			assert.Equal(t, td.fingerprint, td.pubkey.Fingerprint,
				"%s fingerprint %s does not match %s", td.name, td.pubkey.Fingerprint, td.fingerprint)
		})
	}
}

func TestFingerprintLegacyGeneration(t *testing.T) {
	type test struct {
		name        string
		pubkey      *models.Pubkey
		fingerprint string
	}

	tests := []test{
		{"ed25519", factories.NewPubkeyED25519(), "ee:f1:d4:62:99:ab:17:d9:3b:00:66:62:32:b2:55:9e"},
		{"rsa", factories.NewPubkeyRSA(), "89:c5:99:b5:33:48:1c:84:be:da:cb:97:45:b0:4a:ee"},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			err := models.Transform(context.Background(), td.pubkey)
			require.NoError(t, err)
			assert.Equal(t, td.fingerprint, td.pubkey.FingerprintLegacy,
				"%s fingerprint %s does not match %s", td.name, td.pubkey.FingerprintLegacy, td.fingerprint)
		})
	}
}
