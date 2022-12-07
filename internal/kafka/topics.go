package kafka

import "github.com/RHEnVision/provisioning-backend/internal/config"

// topic requests
var (
	availabilityStatusRequestTopicReq = "platform.provisioning.internal.availability-check"
	sendEventsToSourcesTopicReq       = "platform.sources.event-stream"
)

// topics after clowder mapping
var (
	AvailabilityStatusRequestTopic string
	SourcesEventStreamTopic        string
)

// InitializeTopicRequests performs clowder mapping of topics.
func InitializeTopicRequests() {
	AvailabilityStatusRequestTopic = config.TopicName(availabilityStatusRequestTopicReq)
	SourcesEventStreamTopic = config.TopicName(sendEventsToSourcesTopicReq)
}
