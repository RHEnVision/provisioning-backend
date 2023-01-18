package ctxval

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

// Identity returns identity header struct or nil when not set.
func Identity(ctx context.Context) identity.XRHID {
	return identity.Get(ctx)
}

// WithIdentity returns context copy with identity.
func WithIdentity(ctx context.Context, id identity.XRHID) context.Context {
	return context.WithValue(ctx, identity.Key, id)
}

// WithIdentityFrom64 returns context copy with identity parsed from base64-encoded JSON string.
func WithIdentityFrom64(ctx context.Context, id string) (context.Context, error) {
	idRaw, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("unable to b64 decode x-rh-identity %w", err)
	}

	var jsonData identity.XRHID
	err = json.Unmarshal(idRaw, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json %w", err)
	}

	return context.WithValue(ctx, identity.Key, jsonData), nil
}
