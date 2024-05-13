package stubs

import (
	"context"
	"fmt"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
)

type gcpCtxKeyType string

var ipCounter = 10

var (
	gcpCtxKey        gcpCtxKeyType = "gcp-interface"
	serviceGCPCtxKey gcpCtxKeyType = "gcp-service-interface"
)

type (
	GCPClientStub struct {
		Instances []*string
	}
	GCPServiceClientStub struct{}
)

func newGCPCustomerClientStub(ctx context.Context, auth *clients.Authentication) (clients.GCP, error) {
	return getCustomerGCPClientStub(ctx, auth)
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

func getCustomerGCPClientStub(ctx context.Context, auth *clients.Authentication) (*GCPClientStub, error) {
	var si *GCPClientStub
	var err error
	var ok bool
	if si, ok = ctx.Value(gcpCtxKey).(*GCPClientStub); !ok {
		err = ErrContextRead
	}
	return si, err
}

func getServiceGCPClientStub(ctx context.Context) (si clients.ServiceGCP, err error) {
	var ok bool
	if si, ok = ctx.Value(serviceGCPCtxKey).(*GCPServiceClientStub); !ok {
		err = ErrContextRead
	}
	return si, err
}

func CountStubInstancesGCP(ctx context.Context) int {
	client, err := getCustomerGCPClientStub(ctx, &clients.Authentication{})
	if err != nil {
		return 0
	}
	return len(client.Instances)
}

func (mock *GCPClientStub) ListAllRegions(ctx context.Context) ([]clients.Region, error) {
	return nil, nil
}

func (mock *GCPClientStub) Status(ctx context.Context) error {
	return nil
}

func (mock *GCPClientStub) GetInstanceDescriptionByID(ctx context.Context, id, zone string) (*clients.InstanceDescription, error) {
	for _, instanceID := range mock.Instances {
		if ptr.From(instanceID) == id {
			instanceDesc := &clients.InstanceDescription{ID: id, IPv4: fmt.Sprintf("10.0.0.%v", ipCounter)}
			ipCounter = ipCounter + 1
			return instanceDesc, nil
		}
	}
	return nil, ErrMissingInstanceID
}

func (mock *GCPClientStub) ListLaunchTemplates(ctx context.Context) ([]*clients.LaunchTemplate, string, error) {
	return nil, "", nil
}

func (mock *GCPClientStub) InsertInstances(ctx context.Context, params *clients.GCPInstanceParams, amount int64) ([]*string, *string, error) {
	for i := 0; i < int(amount); i++ {
		ID := fmt.Sprintf("300394200587658274%s", strconv.Itoa(len(mock.Instances)+1))
		mock.Instances = append(mock.Instances, &ID)
	}
	ids, err := mock.ListInstancesIDsByLabel(ctx, params.UUID)
	return ids, ptr.To("operation-1686646674436-5fdff07e43209-66146b7e-f3f65ec5"), err
}

func (mock *GCPClientStub) ListInstancesIDsByLabel(ctx context.Context, _ string) ([]*string, error) {
	return mock.Instances, nil
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
