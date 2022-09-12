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

	// GetArn returns ARN associated with provisioning app for given sourceId
	GetArn(ctx context.Context, sourceId ID) (string, error)

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

// GetCustomerEC2Client returns EC2 interface implementation. There are currently
// two implementations available: HTTP and stub
var GetCustomerEC2Client func(ctx context.Context, arn string, region string) (EC2, error)

// Sources interface provides access to the AWS EC2 API
type EC2 interface {
	// ImportPubkey imports new ssh key-pair with given tag returning its AWS ID
	ImportPubkey(key *models.Pubkey, tag string) (string, error)

	// DeleteSSHKey deletes a given ssh key-pair found by AWS ID
	DeleteSSHKey(handle string) error

	// ListInstanceTypesWithPaginator lists available instance types
	ListInstanceTypesWithPaginator() ([]types.InstanceTypeInfo, error)

	// RunInstances launches one or more instances
	RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error)
}

// Caller is responsible for closing the client using Close() call
var GetGCPClient func(ctx context.Context) (GCP, error)

type GCP interface {
	// Close performs close on the gRPC client
	Close()

	// RunInstances launches one or more instances
	RunInstances(ctx context.Context, projectID string, namePattern *string, imageName *string, amount int64, machineType string, zone string, keyBody string) error
}
