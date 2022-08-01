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
