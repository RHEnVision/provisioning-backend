package identity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type Principal = identity.XRHID

// Identity returns identity header struct or nil when not set.
func Identity(ctx context.Context) Principal {
	val := ctx.Value(identity.Key)
	if val == nil {
		return Principal{}
	}
	return val.(Principal)
}

// IdentityHeader returns identity header (base64-encoded JSON)
func IdentityHeader(ctx context.Context) string {
	return identity.GetIdentityHeader(ctx)
}

// WithIdentity returns context copy with identity.
func WithIdentity(ctx context.Context, id Principal) context.Context {
	return context.WithValue(ctx, identity.Key, id)
}

// WithIdentityFrom64 returns context copy with identity parsed from base64-encoded JSON string.
func WithIdentityFrom64(ctx context.Context, id string) (context.Context, error) {
	idRaw, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("unable to b64 decode x-rh-identity %w", err)
	}

	var jsonData Principal
	err = json.Unmarshal(idRaw, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json %w", err)
	}

	return context.WithValue(ctx, identity.Key, jsonData), nil
}
