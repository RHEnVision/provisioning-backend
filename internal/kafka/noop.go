package kafka

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

// Broker that does nothing
type noopBroker struct{}

var _ Broker = &noopBroker{}

func (s *noopBroker) Consume(ctx context.Context, topic string, group string, handler func(ctx context.Context, message *GenericMessage)) {
	logger := ctxval.Logger(ctx)
	logger.Warn().Msg("Consume loop not started (Kafka not configured)")
}

func (s *noopBroker) Send(ctx context.Context, messages ...*GenericMessage) error {
	logger := ctxval.Logger(ctx)
	logger.Warn().Msgf("Throwing away %d messages (Kafka not configured)", len(messages))

	return nil
}
