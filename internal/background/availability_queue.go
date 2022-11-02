package background

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
)

// buffered channel for incoming requests
var kafkaAvailabilityRequest = make(chan *kafka.GenericMessage, availabilityStatusBatchSize)

// EnqueueAvailabilityStatusRequest prepares a status request check to be sent in the next
// batch to the platform kafka. Messages can be delayed up to several seconds until sent.
// The function can block if the enqueueing channel is full.
func EnqueueAvailabilityStatusRequest(msg *kafka.GenericMessage) {
	kafkaAvailabilityRequest <- msg
}

// send a message to the background worker for availability check
func send(ctx context.Context, messages ...*kafka.GenericMessage) {
	err := kafka.Send(ctx, messages...)
	if err != nil {
		ctxval.Logger(ctx).Warn().Err(err).Msg("Unable to send availability check messages")
	}
}

// main sending loop: takes messages enqueued via EnqueueAvailabilityStatusRequest and sends them to the kafka
func sendAvailabilityRequestMessages(ctx context.Context, batchSize int, tickDuration time.Duration) {
	logger := ctxval.Logger(ctx)
	ticker := time.NewTicker(tickDuration)
	messageBuffer := make([]*kafka.GenericMessage, 0, batchSize)

	for {
		select {
		case msg := <-kafkaAvailabilityRequest:
			messageBuffer = append(messageBuffer, msg)
			length := len(messageBuffer)

			if length >= batchSize {
				logger.Trace().Int("messages", length).Msgf("Sent %d availability request messages (full buffer)", length)
				send(ctx, messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		case <-ticker.C:
			length := len(messageBuffer)

			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sent %d availability request messages (tick)", length)
				send(ctx, messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		case <-ctx.Done():
			ticker.Stop()
			length := len(messageBuffer)

			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sent %d availability request messages (cancel)", length)
				send(ctx, messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		}
	}
}
