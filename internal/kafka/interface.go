package kafka

import (
	"context"
)

type Broker interface {
	// Send one or more messages to the kafka
	Send(ctx context.Context, messages ...*GenericMessage) error

	// Consume messages of a single topic in a loop. Blocking call, use context cancellation to stop.
	Consume(ctx context.Context, topic string, group string, handler func(ctx context.Context, message *GenericMessage))
}

var broker Broker = &noopBroker{}

//nolint:wrapcheck
func Send(ctx context.Context, messages ...*GenericMessage) error {
	return broker.Send(ctx, messages...)
}

func Consume(ctx context.Context, topic string, group string, handler func(ctx context.Context, message *GenericMessage)) {
	broker.Consume(ctx, topic, group, handler)
}
