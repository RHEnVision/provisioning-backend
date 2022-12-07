package kafka

import (
	"context"

	"github.com/redhatinsights/platform-go-middlewares/identity"
)

type StatusType string

const (
	StatusUnavailable StatusType = "unavailable"
	StatusAvaliable   StatusType = "available"
)

type SourceResult struct {
	SourceID string `json:"resource_id"`

	// Resource type of the source
	ResourceType string `json:"resource_type"`

	Status StatusType `json:"status"`

	Err error `json:"error"`

	Identity identity.XRHID `json:"-"`
}

func (sr SourceResult) GenericMessage(ctx context.Context) (GenericMessage, error) {
	return genericMessage(ctx, sr, sr.SourceID, SourcesEventStreamTopic)
}
