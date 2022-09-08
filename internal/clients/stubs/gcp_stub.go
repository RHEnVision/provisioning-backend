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

func getGCPClientStub(ctx context.Context, auth *clients.Authentication) (si clients.GCP, err error) {
	var ok bool
	if si, ok = ctx.Value(gcpCtxKey).(*GCPClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *GCPClientStub) Status(ctx context.Context) error {
	return nil
}

func (mock *GCPClientStub) ListAllRegionsAndZones(ctx context.Context) ([]clients.Region, []clients.Zone, error) {
	regions := []clients.Region{
		"us-east1",
		"us-west1",
	}
	zones := []clients.Zone{
		"us-east1-b",
		"us-east1-c",
		"us-east1-d",
		"us-west1-a",
		"us-west1-b",
		"us-west1-c",
	}
	return regions, zones, nil
}

func (mock *GCPClientStub) RunInstances(ctx context.Context, namePattern *string, imageName *string, amount int64, machineType string, zone string, keyBody string) (*string, error) {
	return nil, NotImplementedErr
}
