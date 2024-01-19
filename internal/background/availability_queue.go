package background

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
)

// buffered channel for incoming requests: the length is bigger than the batch size to have room
// for additional messages when first batch is processed.
var kafkaAvailabilityRequest = make(chan *kafka.GenericMessage, 5*availabilityStatusBatchSize)

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

var sendWG sync.WaitGroup

// send a message to the background worker for availability check
func send(ctx context.Context, method string, messages ...*kafka.GenericMessage) {
	// copy buffer into temporary variable
	sendBuffer := make([]*kafka.GenericMessage, len(messages))
	copy(sendBuffer, messages)

	// send on the background
	sendWG.Add(1)
	go func() {
		defer sendWG.Done()
		metrics.ObserveAvailabilitySendDuration(func() {
			logger := zerolog.Ctx(ctx)
			logger.Trace().Int("messages", len(sendBuffer)).
				Msgf("Sending %d availability request messages (%s)", len(sendBuffer), method)

			err := kafka.Send(ctx, sendBuffer...)
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).Msg("Unable to send availability check messages")
			}
		})
	}()
}

// main sending loop: takes messages enqueued via EnqueueAvailabilityStatusRequest and sends them to the kafka
func sendAvailabilityRequestMessages(ctx context.Context, batchSize int, tickDuration time.Duration) {
	ticker := time.NewTicker(tickDuration)
	messageBuffer := make([]*kafka.GenericMessage, 0, batchSize)
	defer sendWG.Wait()

	for {
		select {
		case msg := <-kafkaAvailabilityRequest:
			messageBuffer = append(messageBuffer, msg)
			length := len(messageBuffer)

			if length >= batchSize {
				send(ctx, "full buffer", messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		case <-ticker.C:
			length := len(messageBuffer)

			if length > 0 {
				send(ctx, "tick", messageBuffer...)
				messageBuffer = messageBuffer[:0]
			}
		case <-ctx.Done():
			ticker.Stop()
			length := len(messageBuffer)

			if length > 0 {
				send(ctx, "cancel", messageBuffer...)
			}

			return
		}
	}
}
