package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/identity"
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

func genericMessage(ctx context.Context, m any, key string, topic string) (GenericMessage, error) {
	payload, err := json.Marshal(m)
	if err != nil {
		return GenericMessage{}, fmt.Errorf("unable to marshal message: %w", err)
	}

	id := identity.Identity(ctx)

	return GenericMessage{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		// Keep headers written in lowercase to match sources comparison.
		Headers: GenericHeaders(
			"content-type", "application/json",
			"x-rh-identity", identity.IdentityHeader(ctx),
			"x-rh-sources-org-id", id.Identity.OrgID,
			"x-rh-sources-account-number", id.Identity.AccountNumber,
			"event_type", "availability_status",
		),
	}, nil
}
