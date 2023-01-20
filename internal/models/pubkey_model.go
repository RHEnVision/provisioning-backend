package models

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/ssh"
)

// Pubkey represents SSH public key that can be deployed to clients.
type Pubkey struct {
	// Set to true to skip model validation and transformation during save.
	SkipValidation bool `db:"-" json:"-"`

	// Required auto-generated PK.
	ID int64 `db:"id" json:"id"`

	// Associated Account model. Required.
	AccountID int64 `db:"account_id" json:"-"`

	// User-facing name. Required.
	Name string `db:"name" json:"name" validate:"required"`

	// Public key body encoded in base64 (.pub format). Required.
	Body string `db:"body" json:"body" validate:"required,sshPubkey"`

	// Key type: "ssh-ed25519" or "ssh-rsa".
	Type string `db:"type" json:"type" validate:"omitempty,oneof=test ssh-rsa ssh-ed25519"`

	// SHA256 base64 encoded fingerprint with padding without any prefix. Note OpenSSH
	// typically prints the fingerprint without padding: ssh-keygen -l -f $HOME/.ssh/key.pub
	// The length is exactly 45 characters including the padding.
	// Example: "gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk="
	Fingerprint string `db:"fingerprint" json:"fingerprint" validate:"omitempty,len=44"`

	// MD5 fingerprint stored as hexadecimal with colons without any prefix. To generate
	// such fingerprint: ssh-keygen -l -E md5 -f $HOME/.ssh/key.pub
	// Example: "89:c5:99:b5:33:48:1c:84:be:da:cb:97:45:b0:4a:ee"
	FingerprintLegacy string `db:"fingerprint_legacy" json:"fingerprint_legacy" validate:"omitempty,len=47"`
}

// FindAwsFingerprint returns suitable fingerprint for searching AWS key-pairs.
func (pk *Pubkey) FindAwsFingerprint(ctx context.Context) string {
	switch pk.Type {
	case "ssh-rsa":
		fp, err := ssh.GenerateAWSFingerprint([]byte(pk.Body))
		if err != nil {
			ctxval.Logger(ctx).Warn().Err(err).Msgf("Unable to generate AWS fingerprint for pubkey %s: %s", pk.Name, err.Error())
			return ""
		}
		return string(fp)
	case "ssh-ed25519":
		return pk.Fingerprint
	default:
		return ""

	}
}
