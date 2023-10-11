package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// GenericMessage is a platform independent message.
type GenericMessage struct {
	// Topic of the message. Some producers already have associated topic, in that case Topic from the message will be ignored.
	Topic string

	// Key is used for topic partitioning. Can be nil.
	Key []byte

	// Value is the payload. Typically, a JSON marshaled data.
	Value []byte

	// List of key-value pairs for each message.
	Headers []GenericHeader
}

type GenericHeader struct {
	Key   string
	Value string
}

// GenericHeaders returns slice of headers
func GenericHeaders(args ...string) []GenericHeader {
	if len(args)%2 != 0 {
		panic("generic headers: odd amount of arguments")
	}

	result := make([]GenericHeader, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		gh := GenericHeader{
			Key:   args[i],
			Value: args[i+1],
		}
		result = append(result, gh)
	}

	return result
}

// NativeMessage represents a native (kafka) message. It can be converted to GenericMessage.
type NativeMessage interface {
	// GenericMessage returns a generic message that is platform independent.
	GenericMessage(ctx context.Context) (GenericMessage, error)
}

// NewMessageFromKafka converts generic message to native message
func NewMessageFromKafka(km *kafka.Message) *GenericMessage {
	headers := make([]GenericHeader, len(km.Headers))
	for i, h := range km.Headers {
		gh := GenericHeader{
			Key:   h.Key,
			Value: string(h.Value),
		}
		headers[i] = gh
	}

	return &GenericMessage{
		Topic:   km.Topic,
		Key:     km.Key,
		Value:   km.Value,
		Headers: headers,
	}
}

// KafkaMessage converts from generic to native message.
func (m GenericMessage) KafkaMessage() kafka.Message {
	headers := make([]kafka.Header, len(m.Headers))
	for i, gh := range m.Headers {
		header := kafka.Header{
			Key:   gh.Key,
			Value: []byte(gh.Value),
		}
		headers[i] = header
		i += 1
	}

	return kafka.Message{
		Topic:   m.Topic,
		Key:     m.Key,
		Value:   m.Value,
		Headers: headers,
	}
}

func (m GenericMessage) Header(name string) string {
	for _, h := range m.Headers {
		if strings.EqualFold(h.Key, name) {
			return h.Value
		}
	}
	return ""
}

func genericMessage(ctx context.Context, m any, key string, topic string) (GenericMessage, error) {
	payload, err := json.Marshal(m)
	if err != nil {
		return GenericMessage{}, fmt.Errorf("unable to marshal message: %w", err)
	}

	// This will panic when identity was not present in the context (no error handling possible)
	id := identity.Identity(ctx)
	// Keep headers written in lowercase to match sources comparison.
	headers := GenericHeaders(
		"content-type", "application/json",
		"x-rh-identity", identity.IdentityHeader(ctx),
		"x-rh-sources-org-id", id.Identity.OrgID,
		"x-rh-sources-account-number", id.Identity.AccountNumber,
		"event_type", "availability_status",
	)

	if config.Telemetry.Enabled {
		var traceHeaders propagation.MapCarrier = make(map[string]string)
		otel.GetTextMapPropagator().Inject(ctx, traceHeaders)
		for name, value := range traceHeaders {
			headers = append(headers, GenericHeader{Key: name, Value: value})
		}
	}

	return GenericMessage{
		Topic:   topic,
		Key:     []byte(key),
		Value:   payload,
		Headers: headers,
	}, nil
}

func headersMap(headers []GenericHeader) map[string]string {
	hMap := make(map[string]string, len(headers))
	for _, genericHeader := range headers {
		hMap[genericHeader.Key] = genericHeader.Value
	}
	return hMap
}
