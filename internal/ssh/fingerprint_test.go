package ssh_test

import (
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ssh"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFingerprintPEMGeneration(t *testing.T) {
	type test struct {
		name        string
		pubkey      *models.Pubkey
		fingerprint string
	}

	tests := []test{
		{"ed25519", factories.NewPubkeyED25519(), "e3:8d:76:a4:f2:78:29:f5:6d:0b:95:5c:e9:80:47:85"},
		{"rsa", factories.NewPubkeyRSA(), "c4:ba:72:45:16:a9:2c:39:c3:99:8d:e7:16:01:9c:77"},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			fp, err := ssh.GenerateAWSFingerprint([]byte(td.pubkey.Body))
			assert.NoError(t, err)
			assert.Equal(t, td.fingerprint, string(fp),
				"%s fingerprint %s does not match %s", td.name, string(fp), td.fingerprint)
		})
	}
}

func TestFingerprintUnsupported(t *testing.T) {
	pk := factories.NewPubkeyDSS()
	_, err := ssh.GenerateAWSFingerprint([]byte(pk.Body))
	require.ErrorContains(t, err, "x509: unsupported public key")
}
