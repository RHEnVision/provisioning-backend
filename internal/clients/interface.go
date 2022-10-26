package clients

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// GetSourcesClient returns Sources interface implementation. There are currently
// two implementations available: HTTP and stub
var GetSourcesClient func(ctx context.Context) (Sources, error)

// Sources interface provides access to the Sources backend service API
type Sources interface {
	// ListProvisioningSources returns all sources that have provisioning credentials assigned
	ListProvisioningSources(ctx context.Context) ([]*Source, error)

	// GetArn returns authentication associated with provisioning app for given sourceId
	GetAuthentication(ctx context.Context, sourceId ID) (*Authentication, error)

	// GetProvisioningTypeId returns provisioning type ID
	GetProvisioningTypeId(ctx context.Context) (string, error)

	// Ready returns readiness information
	Ready(ctx context.Context) error
}

// GetImageBuilderClient returns ImageBuilder interface implementation. There are currently
// two implementations available: HTTP and stub
var GetImageBuilderClient func(ctx context.Context) (ImageBuilder, error)

// ImageBuilder interface provides access to the Image Builder backend service API
type ImageBuilder interface {
	// GetAWSAmi returns related AWS image AMI identifier
	GetAWSAmi(ctx context.Context, composeID string) (string, error)

	// GetGCPImageName returns GCP image name
	GetGCPImageName(ctx context.Context, composeID string) (string, error)

	// Ready returns readiness information
	Ready(ctx context.Context) error
}

// ClientStatuser provides a function to test client connection. Since most clouds do not
// provide any "ping" or "status" call, it is usually implemented via some "cheap" operation
// which is fast and returns minimum amount of data (e.g. list regions or ssh-keys).
type ClientStatuser interface {
	Status(ctx context.Context) error
}

// GetEC2Client returns an EC2 facade interface with assumed role.
var GetEC2Client func(ctx context.Context, auth *Authentication, region string) (EC2, error)

// GetServiceEC2Client returns an EC2 client for the service account.
var GetServiceEC2Client func(ctx context.Context, region string) (EC2, error)

type EC2 interface {
	ClientStatuser

	// ListAllRegions returns list of all EC2 regions
	ListAllRegions(ctx context.Context) ([]Region, error)

	// ListAllZones returns list of all EC2 zones within a Region
	ListAllZones(ctx context.Context, region Region) ([]Zone, error)

	// ImportPubkey imports new ssh key-pair with given tag returning its AWS ID
	ImportPubkey(ctx context.Context, key *models.Pubkey, tag string) (string, error)

	// DeleteSSHKey deletes a given ssh key-pair found by AWS ID
	DeleteSSHKey(ctx context.Context, handle string) error

	ListInstanceTypesWithPaginator(ctx context.Context) ([]*InstanceType, error)

	// RunInstances launches one or more instances
	RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error)
}

// GetAzureClient returns an Azure facade interface.
var GetAzureClient func(ctx context.Context, auth *Authentication) (Azure, error)

type Azure interface {
	ClientStatuser
}

// GetGCPClient returns a GCP facade interface of the customer.
var GetGCPClient func(ctx context.Context, auth *Authentication) (GCP, error)

// GetServiceGCPClient returns a GCP client for the service account.
var GetServiceGCPClient func(ctx context.Context) (ServiceGCP, error)

type ServiceGCP interface {
	// RegisterInstanceTypes
	RegisterInstanceTypes(ctx context.Context, instanceTypes *RegisteredInstanceTypes, regionalTypes *RegionalTypeAvailability) error

	// ListMachineTypes returns list of all GCP machine types
	ListMachineTypes(ctx context.Context, zone string) ([]*InstanceType, error)

	// ListAllRegionsAndZones returns list of all GCP regions
	ListAllRegionsAndZones(ctx context.Context) ([]Region, []Zone, error)
}
type GCP interface {
	ClientStatuser

	// ListAllRegions returns list of all GCP regions
	ListAllRegions(ctx context.Context) ([]Region, error)

	// InsertInstances launches one or more instances
	InsertInstances(ctx context.Context, namePattern *string, imageName *string, amount int64, machineType string, zone string, keyBody string) (*string, error)
}
