package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type gcpCtxKeyType string

var gcpCtxKey gcpCtxKeyType = "gcp-interface"

type GCPClientStub struct{}

func init() {
	clients.GetGCPClient = getGCPClientStub
}

// GCPClient
func WithGCPClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, gcpCtxKey, &GCPClientStub{})
	return ctx
}

func getGCPClientStub(ctx context.Context) (si clients.GCP, err error) {
	var ok bool
	if si, ok = ctx.Value(gcpCtxKey).(*GCPClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *GCPClientStub) Close() {
}

func (mock *GCPClientStub) RunInstances(ctx context.Context, projectID string, namePattern *string, imageName *string, amount int64, machineType string, zone string, keyBody string) error {
	return NotImplementedErr
}
