package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/ptr"
)

type ec2CtxKeyType string

var ec2CtxKey ec2CtxKeyType = "ec2-interface"

type EC2ClientStub struct{}

func init() {
	clients.GetCustomerEC2ClientWithRegion = getEC2ClientStubWithRegion
}

// EC2Client
func WithEC2Client(parent context.Context) context.Context {
	ctx := context.WithValue(parent, ec2CtxKey, &EC2ClientStub{})
	return ctx
}

func getEC2ClientStubWithRegion(ctx context.Context, _ string, _ string) (si clients.EC2, err error) {
	var ok bool
	if si, ok = ctx.Value(ec2CtxKey).(*EC2ClientStub); !ok {
		err = &contextReadError{}
	}
	return si, err
}

func (mock *EC2ClientStub) ImportPubkey(key *models.Pubkey, tag string) (string, error) {
	return "", nil
}

func (mock *EC2ClientStub) DeleteSSHKey(handle string) error {
	return nil
}

func (mock *EC2ClientStub) ListInstanceTypesWithPaginator() ([]types.InstanceTypeInfo, error) {
	return []types.InstanceTypeInfo{
		{
			InstanceType: types.InstanceTypeA12xlarge,
			VCpuInfo: &types.VCpuInfo{
				DefaultCores: ptr.Int32(2),
				DefaultVCpus: ptr.Int32(2),
			},
			MemoryInfo: &types.MemoryInfo{
				SizeInMiB: ptr.Int64(22),
			},
			ProcessorInfo: &types.ProcessorInfo{
				SupportedArchitectures: []types.ArchitectureType{types.ArchitectureTypeX8664, types.ArchitectureTypeArm64},
			},
		},
		{
			InstanceType: types.InstanceTypeC5Xlarge,
			VCpuInfo: &types.VCpuInfo{
				DefaultCores: ptr.Int32(2),
				DefaultVCpus: ptr.Int32(2),
			},
			MemoryInfo: &types.MemoryInfo{
				SizeInMiB: ptr.Int64(22),
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
