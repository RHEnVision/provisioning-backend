package models

import (
	"crypto/md5" //#nosec
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Pubkey represents SSH public key that can be deployed to clients.
type Pubkey struct {
	// Required auto-generated PK.
	ID int64 `db:"id" json:"id"`

	// Associated Account model. Required.
	AccountID int64 `db:"account_id" json:"-"`

	// User-facing name. Required.
	Name string `db:"name" json:"name" validate:"required"`

	// Public key body encoded in base64 (.pub format). Required.
	Body string `db:"body" json:"body" validate:"required,sshPubkey"`

	// Public key SHA256 fingerprint. Generated (read-only).
	Fingerprint string `db:"fingerprint" json:"fingerprint"`
}

// FingerprintAWS calculates fingerprint used by AWS.
// AWS calculates MD5 sums differently then all the other tools and use format DER as the base for the sub.
// see https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/verify-keys.html for more details.
func (pk *Pubkey) FingerprintAWS() (string, error) {
	pkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pk.Body))
	if err != nil {
		return "", fmt.Errorf("failed to calculate fingerprint for %s: %w", pk.Name, err)
	}
	// when rsa key
	if pkey.Type() == "ssh-rsa" {
		return fingerprintRsaAWS(pkey)
	}
	// when ed25519 key
	// this is the same as what we store, but better be sure here
	return strings.TrimLeft(ssh.FingerprintSHA256(pkey), "SHA256:") + "=", nil
}

// fingerprintRsaAWS calculates fingerprint for rsa keys
func fingerprintRsaAWS(key ssh.PublicKey) (string, error) {
	parsedCryptoKey := key.(ssh.CryptoPublicKey)

	// Finally, we can convert back to an *rsa.PublicKey
	pub := parsedCryptoKey.CryptoPublicKey().(*rsa.PublicKey)

	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", fmt.Errorf("failed to calculate rsa key fingerprint: %w", err)
	}
	md5sum := md5.Sum(der) //#nosec
	hexarray := make([]string, len(md5sum))
	for i, c := range md5sum {
		hexarray[i] = hex.EncodeToString([]byte{c})
	}
	return strings.Join(hexarray, ":"), nil
}
