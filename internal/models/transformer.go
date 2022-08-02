package models

import (
	"context"
	"fmt"
	"reflect"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/go-playground/mold/v4"
	"golang.org/x/crypto/ssh"
)

var transform *mold.Transformer

func init() {
	transform = mold.New()
	transform.RegisterStructLevel(generateFingerprint, Pubkey{})
}

func Transform(ctx context.Context, model interface{}) error {
	err := transform.Struct(ctx, model)
	if err != nil {
		return fmt.Errorf("transformation error: %w", err)
	}
	return nil
}

func generateFingerprint(ctx context.Context, sl mold.StructLevel) error {
	logger := ctxval.Logger(ctx)
	pk := sl.Struct().Interface().(Pubkey)
	pkey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pk.Body))
	if err != nil {
		return fmt.Errorf("unable to calculate fingerprint for %s: %w", pk.Name, err)
	}
	pk.Fingerprint = ssh.FingerprintSHA256(pkey)
	sl.Struct().Set(reflect.ValueOf(pk))
	logger.Trace().Msgf("Calculated SSH fingerprint for %s: %s", pk.Name, pk.Fingerprint)
	return nil
}
