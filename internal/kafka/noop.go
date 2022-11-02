package kafka

import (
	"context"
)

// Broker that does nothing
type noopBroker struct{}

var _ Broker = &noopBroker{}

func (s *noopBroker) Consume(ctx context.Context, topic string, handler func(ctx context.Context, message *GenericMessage)) {
}

func (s *noopBroker) Send(_ context.Context, messages ...*GenericMessage) error {
	return nil
}
