package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ec2CtxKeyType string

var ec2CtxKey ec2CtxKeyType = "ec2-interface"

type EC2ClientStub struct{}

func init() {
	clients.GetCustomerEC2Client = getEC2ClientStubWithRegion
}

// EC2Client
func WithEC2Client(parent context.Context) context.Context {
	ctx := context.WithValue(parent, ec2CtxKey, &EC2ClientStub{})
	return ctx
}

func getEC2ClientStubWithRegion(ctx context.Context, _ *clients.Authentication, _ string) (si clients.EC2, err error) {
	var ok bool
	if si, ok = ctx.Value(ec2CtxKey).(*EC2ClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *EC2ClientStub) Status(ctx context.Context) error {
	return nil
}

func (mock *EC2ClientStub) ImportPubkey(ctx context.Context, key *models.Pubkey, tag string) (string, error) {
	return "", nil
}

func (mock *EC2ClientStub) DeleteSSHKey(ctx context.Context, handle string) error {
	return nil
}

func (mock *EC2ClientStub) ListAllRegions(ctx context.Context) ([]clients.Region, error) {
	return []clients.Region{
		"us-east-1",
		"eu-central-1",
	}, nil
}

func (mock *EC2ClientStub) ListAllZones(ctx context.Context, region clients.Region) ([]clients.Zone, error) {
	return []clients.Zone{
		"us-east-1a",
		"us-east-1b",
		"us-east-1c",
		"eu-central-1a",
		"eu-central-1b",
		"eu-central-1c",
	}, nil
}

func (mock *EC2ClientStub) ListInstanceTypesWithPaginator(ctx context.Context) ([]types.InstanceTypeInfo, error) {
	return []types.InstanceTypeInfo{
		{
			InstanceType: types.InstanceTypeA12xlarge,
			VCpuInfo: &types.VCpuInfo{
				DefaultCores: ptr.ToInt32(2),
				DefaultVCpus: ptr.ToInt32(2),
			},
			MemoryInfo: &types.MemoryInfo{
				SizeInMiB: ptr.ToInt64(22),
			},
			ProcessorInfo: &types.ProcessorInfo{
				SupportedArchitectures: []types.ArchitectureType{types.ArchitectureTypeX8664, types.ArchitectureTypeArm64},
			},
		},
		{
			InstanceType: types.InstanceTypeC5Xlarge,
			VCpuInfo: &types.VCpuInfo{
				DefaultCores: ptr.ToInt32(2),
				DefaultVCpus: ptr.ToInt32(2),
			},
			MemoryInfo: &types.MemoryInfo{
				SizeInMiB: ptr.ToInt64(22),
			},
			ProcessorInfo: &types.ProcessorInfo{
				SupportedArchitectures: []types.ArchitectureType{types.ArchitectureTypeX8664},
			},
		},
	}, nil
}

func (mock *EC2ClientStub) RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error) {
	return nil, nil, nil
}
