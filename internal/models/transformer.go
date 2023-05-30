package models

import (
	"context"
	"fmt"
	"reflect"

	"github.com/RHEnVision/provisioning-backend/internal/ssh"
	"github.com/go-playground/mold/v4"
	"github.com/rs/zerolog"
)

var transform *mold.Transformer

func init() {
	transform = mold.New()
	transform.RegisterStructLevel(generateFingerprints, Pubkey{})
}

func Transform(ctx context.Context, model interface{}) error {
	err := transform.Struct(ctx, model)
	if err != nil {
		return fmt.Errorf("transformation error: %w", err)
	}
	return nil
}

// generates fingerprint fields or returns an error for unsupported keys: "x509: unsupported public key"
func generateFingerprints(ctx context.Context, sl mold.StructLevel) error {
	pk := sl.Struct().Interface().(Pubkey)

	pkf, err := ssh.GenerateOpenSSHFingerprints([]byte(pk.Body))
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Str("pubkey", pk.Body).Msg("OpenSSH fingerprint generation error")
		return fmt.Errorf("key generate error %s: %w", pk.Name, err)
	}

	pk.Type = pkf.Type
	pk.Fingerprint = pkf.SHA256
	pk.FingerprintLegacy = pkf.MD5
	sl.Struct().Set(reflect.ValueOf(pk))

	err = validateAWS(ctx, sl)
	if err != nil {
		return fmt.Errorf("key error %s: %w", pk.Name, err)
	}

	return nil
}

// validateAWS tries to generate AWS PEM key during key save because the fingerprint is generated
// on the fly and can fail later.
func validateAWS(ctx context.Context, sl mold.StructLevel) error {
	pk := sl.Struct().Interface().(Pubkey)

	_, err := ssh.GenerateAWSFingerprint([]byte(pk.Body))
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Str("pubkey", pk.Body).Msg("AWS fingerprint validation error")
		return fmt.Errorf("invalid public key type (only ed25519 and rsa keys are supported): %w", err)
	}
	sl.Struct().Set(reflect.ValueOf(pk))

	return nil
}
