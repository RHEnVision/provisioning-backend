package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type ctxKeyType int

const azureCtxKey ctxKeyType = iota

func WithAzureClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, azureCtxKey, &AzureClientStub{})
	return ctx
}

func getAzureClient(ctx context.Context, auth *clients.Authentication) (clients.Azure, error) {
	return getAzureClientStub(ctx)
}

func getAzureClientStub(ctx context.Context) (*AzureClientStub, error) {
	var si *AzureClientStub
	var ok bool
	if si, ok = ctx.Value(azureCtxKey).(*AzureClientStub); !ok {
		return nil, ErrContextRead
	}
	return si, nil
}
