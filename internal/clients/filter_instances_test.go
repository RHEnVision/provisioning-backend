package clients_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
)

func TestEC2MapArchitectures(t *testing.T) {
	result, err := clients.MapArchitectures(context.Background(), "x86-64")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeX8664, result)
	result, err = clients.MapArchitectures(context.Background(), "i386")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeI386, result)
	result, err = clients.MapArchitectures(context.Background(), "arm")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeArm64, result)
	result, err = clients.MapArchitectures(context.Background(), "86_64_mac")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeAppleX8664, result)
	result, err = clients.MapArchitectures(context.Background(), "arm_mac")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeAppleArm64, result)
	_, err = clients.MapArchitectures(context.Background(), "ppc64")
	assert.NotNil(t, err)
}

func TestAzureMapArchitectures(t *testing.T) {
	result, err := clients.MapArchitectures(context.Background(), "x64")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeX8664, result)
	result, err = clients.MapArchitectures(context.Background(), "arm64")
	assert.Nil(t, err)
	assert.Equal(t, clients.ArchitectureTypeArm64, result)
}

func TestNewInstanceTypes(t *testing.T) {
	AWSInstanceTypes := []types.InstanceTypeInfo{
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
	}
	res, err := ec2.NewInstanceTypes(context.Background(), AWSInstanceTypes)
	assert.Nil(t, err)
	// Check that two instance types were created, one per architecture
	assert.Equal(t, len(*res), 3)
	assert.Equal(t, clients.InstanceTypeName("a1.2xlarge"), (*res)[0].Name)
	assert.Equal(t, clients.InstanceTypeName("a1.2xlarge"), (*res)[1].Name)
	assert.Equal(t, clients.InstanceTypeName("c5.xlarge"), (*res)[2].Name)
	// Check that instance types which does not appear in supported_instance_yml are marked as unsupported
	assert.Equal(t, (*res)[0].Supported, false)
	assert.Equal(t, (*res)[1].Supported, false)
	// Check that instance types which appear in supported_instance_yml are marked as supported
	assert.Equal(t, (*res)[2].Supported, true)
}
