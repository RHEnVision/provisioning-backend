package kafka

import (
	"context"
	"encoding/json"
	"fmt"
)

type AvailabilityStatusMessage struct {
	SourceID string `json:"source_id"`
}

func NewAvailabilityStatusMessage(msg *GenericMessage) (*AvailabilityStatusMessage, error) {
	asm := AvailabilityStatusMessage{}
	err := json.Unmarshal(msg.Value, &asm)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal message: %w", err)
	}

	return &asm, nil
}

func (m AvailabilityStatusMessage) GenericMessage(ctx context.Context) (GenericMessage, error) {
	return genericMessage(ctx, m, m.SourceID, AvailabilityStatusRequestTopic)
}
