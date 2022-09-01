package clients

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

var GetSourcesClient func(ctx context.Context) (Sources, error)

type Sources interface {
	// ListProvisioningSources returns all sources that have provisioning credentials assigned
	ListProvisioningSources(ctx context.Context) (*[]Source, error)
	// GetArn returns ARN associated with provisioning app for given sourceId
	GetArn(ctx context.Context, sourceId ID) (string, error)
	// GetProvisioningTypeId might not need exposing
	GetProvisioningTypeId(ctx context.Context) (string, error)
	// Ready returns readiness information
	Ready(ctx context.Context) error
}

var GetImageBuilderClient func(ctx context.Context) (ImageBuilder, error)

type ImageBuilder interface {
	// GetAWSAmi returns related AWS image AMI identifer
	GetAWSAmi(ctx context.Context, composeID string) (string, error)
	// Ready returns readiness information
	Ready(ctx context.Context) error
}

var GetCustomerEC2ClientWithRegion func(ctx context.Context, arn string, region string) (EC2, error)

func GetCustomerEC2Client(ctx context.Context, arn string) (EC2, error) {
	return GetCustomerEC2ClientWithRegion(ctx, arn, config.AWS.Region)
}

type EC2 interface {
	ImportPubkey(key *models.Pubkey, tag string) (string, error)
	DeleteSSHKey(handle string) error
	ListInstanceTypesWithPaginator() ([]types.InstanceTypeInfo, error)
	RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error)
}
