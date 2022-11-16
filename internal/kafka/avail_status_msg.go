package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/redhatinsights/platform-go-middlewares/identity"
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

	id := ctxval.Identity(ctx)

	return GenericMessage{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"x-RH-Identity":               identity.GetIdentityHeader(ctx),
			"X-RH-Sources-Org-Id":         id.Identity.OrgID,
			"X-RH-Sources-Account-Number": id.Identity.AccountNumber,
		},
	}, nil
}
