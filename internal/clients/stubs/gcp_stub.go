package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type gcpCtxKeyType string

var (
	gcpCtxKey        gcpCtxKeyType = "gcp-interface"
	serviceGCPCtxKey gcpCtxKeyType = "gcp-service-interface"
)

type (
	GCPClientStub        struct{}
	GCPServiceClientStub struct{}
)

func init() {
	clients.GetGCPClient = getCustomerGCPClientStub
	clients.GetServiceGCPClient = getServiceGCPClientStub
}

// GCPServiceClient
func WithGCPCServiceClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, serviceGCPCtxKey, &GCPServiceClientStub{})
	return ctx
}

// GCPCustomerClient
func WithGCPCCustomerClient(parent context.Context) context.Context {
	ctx := context.WithValue(parent, gcpCtxKey, &GCPClientStub{})
	return ctx
}

func getCustomerGCPClientStub(ctx context.Context, auth *clients.Authentication) (si clients.GCP, err error) {
	var ok bool
	if si, ok = ctx.Value(gcpCtxKey).(*GCPClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func getServiceGCPClientStub(ctx context.Context) (si clients.ServiceGCP, err error) {
	var ok bool
	if si, ok = ctx.Value(serviceGCPCtxKey).(*GCPServiceClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *GCPClientStub) ListAllRegions(ctx context.Context) ([]clients.Region, error) {
	return nil, nil
}

func (mock *GCPClientStub) Status(ctx context.Context) error {
	return nil
}

func (mock *GCPClientStub) InsertInstances(ctx context.Context, namePattern *string, imageName *string, amount int64, machineType string, zone string, keyBody string) (*string, error) {
	return nil, NotImplementedErr
}

func (mock *GCPServiceClientStub) ListMachineTypes(ctx context.Context, zone string) ([]*clients.InstanceType, error) {
	return nil, nil
}

func (mock *GCPServiceClientStub) RegisterInstanceTypes(ctx context.Context, instanceTypes *clients.RegisteredInstanceTypes, regionalTypes *clients.RegionalTypeAvailability) error {
	return nil
}

func (mock *GCPServiceClientStub) ListAllRegionsAndZones(ctx context.Context) ([]clients.Region, []clients.Zone, error) {
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
