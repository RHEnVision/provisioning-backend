package kafka

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// Broker that does nothing
type noopBroker struct{}

var _ Broker = &noopBroker{}

func (s *noopBroker) Consume(ctx context.Context, topic string, since time.Time, handler func(ctx context.Context, message *GenericMessage)) {
	logger := zerolog.Ctx(ctx)
	logger.Warn().Msg("Consume loop not started (Kafka not configured)")
}

func (s *noopBroker) Send(ctx context.Context, messages ...*GenericMessage) error {
	logger := zerolog.Ctx(ctx)
	logger.Warn().Msgf("Throwing away %d messages (Kafka not configured)", len(messages))

	return nil
}
