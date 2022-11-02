package kafka

import (
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

func (m AvailabilityStatusMessage) GenericMessage() (GenericMessage, error) {
	payload, err := json.Marshal(m)
	if err != nil {
		return GenericMessage{}, fmt.Errorf("unable to marshal message: %w", err)
	}

	return GenericMessage{
		Topic: AvailabilityStatusRequestTopic,
		Key:   []byte(m.SourceID),
		Value: payload,
	}, nil
}
