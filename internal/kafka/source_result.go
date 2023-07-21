package kafka

import (
	"context"
)

type StatusType string

const (
	StatusUnavailable StatusType = "unavailable"
	StatusAvailable   StatusType = "available"
)

type SourceResult struct {
	MessageContext     context.Context `json:"-"` // Carries logger and identity
	ResourceID         string          `json:"resource_id"`
	ResourceType       string          `json:"resource_type"`
	Status             StatusType      `json:"status"`
	Err                error           `json:"-"` // Sources do not support error field
	MissingPermissions []string        `json:"-"` // Sources do not support reason field
}

func (sr SourceResult) GenericMessage(ctx context.Context) (GenericMessage, error) {
	return genericMessage(ctx, sr, sr.ResourceID, SourcesStatusTopic)
}

func (st StatusType) String() string {
	return string(st)
}
