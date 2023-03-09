package clients

import (
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type GCPInstanceParams struct {
	// The string pattern used for the name of the VM.
	NamePattern *string

	// Image Name the instance will be launched from
	ImageName string

	// InstanceType to launch
	MachineType string

	// Zone - to deploy into
	Zone string

	// Pubkey to use for the instance access
	KeyBody string
}

type AWSInstanceParams struct {
	// The template id to use in order to launch an instance
	LaunchTemplateID string

	// ami of the instance will be launched from
	AMI string

	// InstanceType to launch
	InstanceType types.InstanceType

	// Zone - to deploy into
	Zone string

	// Pubkey to use for the instance access
	KeyName string

	// UserData for the instance launch
	UserData []byte
}

// AzureInstanceParams define parameters for a single instance launch on Azure.
type AzureInstanceParams struct {
	// Location - to deploy into
	Location string

	// ResourceGroupName to launch the instance in
	ResourceGroupName string

	// ImageName - the imageID will be inferred as /subscriptions/{subscriptionID}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}
	ImageName string

	// Pubkey to use for the instance access
	Pubkey *models.Pubkey

	// InstanceType to launch
	InstanceType InstanceTypeName

	// UserData for the instance launch
	UserData []byte
}
