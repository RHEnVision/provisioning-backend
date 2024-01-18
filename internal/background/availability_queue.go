package background

import (
	"context"
	"fmt"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
)

// buffered channel for incoming requests
// the length is twice the batch size to have room for additional messages when first batch is processed.
var kafkaAvailabilityRequest = make(chan *kafka.GenericMessage, 2*availabilityStatusBatchSize)

// EnqueueAvailabilityStatusRequest prepares a status request check to be sent in the next
// batch to the platform kafka. Messages can be delayed up to several seconds until sent.
// The function can block if the enqueueing channel is full.
func EnqueueAvailabilityStatusRequest(ctx context.Context, asm *kafka.AvailabilityStatusMessage) error {
	ctx, span := telemetry.StartSpan(ctx, "EnqueueAvailabilityStatusRequest")
	defer span.End()

	msg, err := asm.GenericMessage(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to generate Kafka message")
		return fmt.Errorf("cannot create message: %w", err)
	}

	kafkaAvailabilityRequest <- &msg
	zerolog.Ctx(ctx).Trace().Str("source_id", asm.SourceID).Msgf("Enqueued source id %s availability check", asm.SourceID)
	return nil
}

// send a message to the background worker for availability check
func send(ctx context.Context, messages ...*kafka.GenericMessage) {
	err := kafka.Send(ctx, messages...)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Unable to send availability check messages")
	}
}

// main sending loop: takes messages enqueued via EnqueueAvailabilityStatusRequest and sends them to the kafka
func sendAvailabilityRequestMessages(ctx context.Context, batchSize int, tickDuration time.Duration) {
	logger := zerolog.Ctx(ctx)
	ticker := time.NewTicker(tickDuration)
	messageBuffer := make([]*kafka.GenericMessage, 0, batchSize)

	for {
		select {
		case msg := <-kafkaAvailabilityRequest:
			messageBuffer = append(messageBuffer, msg)
			length := len(messageBuffer)

			if length >= batchSize {
				logger.Trace().Int("messages", length).Msgf("Sending %d availability request messages (full buffer)", length)
				send(ctx, messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		case <-ticker.C:
			length := len(messageBuffer)

			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d availability request messages (tick)", length)
				send(ctx, messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		case <-ctx.Done():
			ticker.Stop()
			length := len(messageBuffer)

			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d availability request messages (cancel)", length)
				send(ctx, messageBuffer...)
			}

			return
		}
	}
}
