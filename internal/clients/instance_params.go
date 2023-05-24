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

	// UUID for instance that was created in a reservation
	UUID string

	// The template name to use in order to launch an instance
	LaunchTemplateName string

	// Zone - to deploy into
	Zone string

	// Pubkey to use for the instance access
	KeyBody string

	// StartupScript contains metadata startup script (GCP tools must be installed on the image)
	StartupScript string
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

	// ImageID - the Image ID in format of full Azure ID as
	// for example /subscriptions/{subscriptionID}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/images/{imageName}
	ImageID string

	// Pubkey to use for the instance access
	Pubkey *models.Pubkey

	// InstanceType to launch
	InstanceType InstanceTypeName

	// UserData for the instance launch
	UserData []byte
}
