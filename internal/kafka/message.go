package kafka

import (
	"context"
	"strings"

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
