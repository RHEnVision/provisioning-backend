package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/google/uuid"
)

const (
	application                  = "image-builder"
	bundle                       = "rhel"
	notificationMessageVersion   = "v2.0.0"
	NotificationSuccessEventType = "launch-success"
	NotificationFailureEventType = "launch-failed"
)

type NotificationEvent struct {
	Payload json.RawMessage `json:"payload"`
}

type NotificationError struct {
	Error string `json:"error"`
}

type notificationRecipients struct {
	OnlyAdmins            bool     `json:"only_admins"`
	IgnoreUserPreferences bool     `json:"ignore_user_preferences"`
	Users                 []string `json:"users"`
}

type NotificationMessage struct {
	Version     string                   `json:"version"`
	Bundle      string                   `json:"bundle"`
	Application string                   `json:"application"`
	EventType   string                   `json:"event_type"`
	Timestamp   string                   `json:"timestamp"`
	AccountID   string                   `json:"account_id"`
	OrgId       string                   `json:"org_id"`
	Context     interface{}              `json:"context"`
	Events      []NotificationEvent      `json:"events"`
	Recipients  []notificationRecipients `json:"recipients"`
	ID          string                   `json:"id"`
}

func (m NotificationMessage) GenericMessage(ctx context.Context) (GenericMessage, error) {
	id := identity.Identity(ctx)

	m.Application = application
	m.Bundle = bundle
	m.Version = notificationMessageVersion
	m.OrgId = id.Identity.OrgID
	m.AccountID = id.Identity.AccountNumber
	m.Timestamp = time.Now().Format("2006-01-02T15:04:05.000") // ISO 8601 REQUIRED
	m.ID = uuid.New().String()
	m.Recipients = []notificationRecipients{}

	payload, err := json.Marshal(m)
	if err != nil {
		return GenericMessage{}, fmt.Errorf("unable to marshal notification message: %w", err)
	}

	return GenericMessage{
		Topic: NotificationTopic,
		Key:   nil,
		Value: payload,
		Headers: GenericHeaders(
			"rh-message-id", m.ID,
			"x-rh-identity", identity.IdentityHeader(ctx),
		),
	}, nil
}
