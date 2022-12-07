package kafka

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/config"
)

// topic requests
var (
	availabilityStatusRequestTopicReq = "platform.provisioning.internal.availability-check"
	sendStatusToSourcesTopicReq       = "platform.sources.status"
)

// topics after clowder mapping
var (
	AvailabilityStatusRequestTopic string
	SourcesStatusTopic             string
)

// InitializeTopicRequests performs clowder mapping of topics.
func InitializeTopicRequests(ctx context.Context) {
	AvailabilityStatusRequestTopic = config.TopicName(ctx, availabilityStatusRequestTopicReq)
	SourcesStatusTopic = config.TopicName(ctx, sendStatusToSourcesTopicReq)
}
