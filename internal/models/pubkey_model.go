package models

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/ssh"
	"github.com/rs/zerolog"
)

var ErrInvalidPubkeyFormat = errors.New("invalid public key format")

// Pubkey represents SSH public key that can be deployed to clients.
type Pubkey struct {
	// Set to true to skip model validation and transformation during save.
	SkipValidation bool `db:"-"`

	// Required auto-generated PK.
	ID int64 `db:"id"`

	// Associated Account model. Required.
	AccountID int64 `db:"account_id"`

	// User-facing name. Required.
	Name string `db:"name" validate:"required"`

	// Public key body encoded in base64 (.pub format). Required.
	Body string `db:"body" validate:"required,sshPubkey"`

	// Key type: "ssh-ed25519" or "ssh-rsa".
	Type string `db:"type" validate:"omitempty,oneof=test ssh-rsa ssh-ed25519"`

	// SHA256 base64 encoded fingerprint with padding without any prefix. Note OpenSSH
	// typically prints the fingerprint without padding: ssh-keygen -l -f $HOME/.ssh/key.pub
	// The length is exactly 45 characters including the padding.
	// Example: "gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk="
	Fingerprint string `db:"fingerprint" validate:"omitempty,len=44"`

	// MD5 fingerprint stored as hexadecimal with colons without any prefix. To generate
	// such fingerprint: ssh-keygen -l -E md5 -f $HOME/.ssh/key.pub
	// Example: "89:c5:99:b5:33:48:1c:84:be:da:cb:97:45:b0:4a:ee"
	FingerprintLegacy string `db:"fingerprint_legacy" validate:"omitempty,len=47"`
}

// FindAwsFingerprint returns suitable fingerprint for searching AWS key-pairs.
func (pk *Pubkey) FindAwsFingerprint(ctx context.Context) string {
	switch pk.Type {
	case "ssh-rsa":
		fp, err := ssh.GenerateAWSFingerprint([]byte(pk.Body))
		if err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Msg("Unable to generate AWS fingerprint for pubkey")
			return ""
		}
		return string(fp)
	case "ssh-ed25519":
		return pk.Fingerprint
	default:
		return ""

	}
}

func (pk *Pubkey) BodyWithUsername(ctx context.Context) (string, error) {
	parts := strings.Split(pk.Body, " ")
	if len(parts) < 2 {
		return "", ErrInvalidPubkeyFormat
	}
	return fmt.Sprintf("gcp-user:%s %s", parts[0], parts[1]), nil
}
