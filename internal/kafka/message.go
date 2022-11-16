package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
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
	Headers map[string]string
}

// NativeMessage represents a native (kafka) message. It can be converted to GenericMessage.
type NativeMessage interface {
	// GenericMessage returns a generic message that is platform independent.
	GenericMessage(ctx context.Context) (GenericMessage, error)
}

// NewMessageFromKafka converts generic message to native message
func NewMessageFromKafka(km *kafka.Message) *GenericMessage {
	headers := make(map[string]string, len(km.Headers))
	for _, h := range km.Headers {
		headers[h.Key] = string(h.Value)
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
	i := 0
	for k, v := range m.Headers {
		header := kafka.Header{
			Key:   k,
			Value: []byte(v),
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
