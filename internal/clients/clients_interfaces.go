package clients

import (
	"context"
)

var GetSourcesClient func(ctx context.Context) (Sources, error)

type Sources interface {
	// ListProvisioningSources returns all sources that have provisioning credentials assigned
	ListProvisioningSources(ctx context.Context) (*[]Source, error)
	// GetArn returns ARN associated with provisioning app for given sourceId
	GetArn(ctx context.Context, sourceId ID) (string, error)
	// GetProvisioningTypeId might not need exposing
	GetProvisioningTypeId(ctx context.Context) (string, error)
}

var GetImageBuilderClient func(ctx context.Context) (ImageBuilder, error)

type ImageBuilder interface {
	// GetAWSAmi returns related AWS image AMI identifer
	GetAWSAmi(ctx context.Context, composeID string) (string, error)
}
