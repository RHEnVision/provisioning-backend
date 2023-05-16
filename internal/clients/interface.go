package clients

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/models"
)

// GetSourcesClient returns Sources interface implementation. There are currently
// two implementations available: HTTP and stub
var GetSourcesClient func(ctx context.Context) (Sources, error)

// Sources interface provides access to the Sources backend service API
type Sources interface {
	// ListProvisioningSourcesByProvider returns sources filtered by provider that have provisioning credentials assigned
	ListProvisioningSourcesByProvider(ctx context.Context, provider models.ProviderType) ([]*Source, error)

	// ListAllProvisioningSources returns all sources that have provisioning credentials assigned
	ListAllProvisioningSources(ctx context.Context) ([]*Source, error)

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

	// GetAzureImageID returns partial image id, that is missing the subscription prefix
	// Full name is /subscriptions/<subscription-id>/resourceGroups/<Group>/providers/Microsoft.Compute/images/<ImageName>
	// GetAzureImageID returns /resourceGroups/<Group>/providers/Microsoft.Compute/images/<ImageName>
	GetAzureImageID(ctx context.Context, composeID string) (string, error)

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

	// ListAllRegions returns list of all EC2 regions.
	ListAllRegions(ctx context.Context) ([]Region, error)

	// ListAllZones returns list of all EC2 zones within a Region.
	ListAllZones(ctx context.Context, region Region) ([]Zone, error)

	// ImportPubkey imports new ssh key-pair with given tag returning its AWS ID.
	ImportPubkey(ctx context.Context, key *models.Pubkey, tag string) (string, error)

	// GetPubkeyName fetches the AWS key name using given pubkey fingerprint.
	GetPubkeyName(ctx context.Context, fingerprint string) (string, error)

	// DeleteSSHKey deletes a given ssh key-pair found by AWS ID.
	DeleteSSHKey(ctx context.Context, handle string) error

	// ListInstanceTypesWithPaginator lists all instance types.
	ListInstanceTypes(ctx context.Context) ([]*InstanceType, error)

	// ListLaunchTemplates lists all launch templates.
	ListLaunchTemplates(ctx context.Context) ([]*LaunchTemplate, error)

	// RunInstances launches one or more instances.
	//
	// All arguments are required except: launchTemplateID (empty string means no template in use).
	//
	RunInstances(ctx context.Context, details *AWSInstanceParams, amount int32, name *string) ([]*string, *string, error)

	// GetAccountId returns AWS account number.
	GetAccountId(ctx context.Context) (string, error)

	CheckPermission(ctx context.Context, auth *Authentication) ([]string, error)

	DescribeInstanceDetails(ctx context.Context, InstanceIds []string) ([]*InstanceDescription, error)
}

// GetAzureClient returns an Azure client with customer's subscription ID.
var GetAzureClient func(ctx context.Context, auth *Authentication) (Azure, error)

// GetServiceAzureClient returns an Azure client for the service account itself.
var GetServiceAzureClient func(ctx context.Context) (ServiceAzure, error)

type Azure interface {
	ClientStatuser

	// TenantId returns current subscription's tenant
	TenantId(ctx context.Context) (AzureTenantId, error)

	// EnsureResourceGroup makes sure that group with give name exists in a location
	EnsureResourceGroup(ctx context.Context, name string, location string) (*string, error)

	// CreateVMs creates multiple Azure virtual machines
	// Returns array of instance IDs and error if something went wrong
	CreateVMs(ctx context.Context, instanceParams AzureInstanceParams, amount int64, vmNamePrefix string) (vmIds []InstanceDescription, err error)

	ListResourceGroups(ctx context.Context) ([]string, error)
}

type ServiceAzure interface {
	RegisterInstanceTypes(ctx context.Context, instanceTypes *RegisteredInstanceTypes, regionalTypes *RegionalTypeAvailability) error
}

// GetGCPClient returns a GCP facade interface.
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

	// InsertInstances launches one or more instances and returns a list of instances ids that were created, the GCP operation name and error
	InsertInstances(ctx context.Context, params *GCPInstanceParams, amount int64) ([]*string, *string, error)
}
