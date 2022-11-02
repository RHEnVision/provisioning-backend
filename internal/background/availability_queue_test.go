package background

import (
	"context"
	"sync"
	"testing"
	"time"

	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"

	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/stretchr/testify/require"
)

func TestQueueNormalSend(t *testing.T) {
	ctx := context.Background()
	_ = kafka.InitializeStubBroker(16)

	wg := sync.WaitGroup{}
	wg.Add(2)
	cct, cancel := context.WithCancel(ctx)
	defer cancel()
	go sendAvailabilityRequestMessages(cct, 8, 10*time.Millisecond)
	go kafka.Consume(cct, kafka.AvailabilityStatusRequestTopic, func(ctx context.Context, msg *kafka.GenericMessage) {
		asm, _ := kafka.NewAvailabilityStatusMessage(msg)
		require.EqualValues(t, "1", asm.SourceID)
		wg.Done()
	})

	msg, _ := kafka.AvailabilityStatusMessage{SourceID: "1"}.GenericMessage()
	EnqueueAvailabilityStatusRequest(&msg)
	EnqueueAvailabilityStatusRequest(&msg)
	wg.Wait()
}

func TestFullQueueSend(t *testing.T) {
	ctx := context.Background()
	_ = kafka.InitializeStubBroker(16)

	wg := sync.WaitGroup{}
	wg.Add(2)
	consumeCtx, consumeCancel := context.WithCancel(ctx)
	senderCtx, senderCancel := context.WithCancel(ctx)
	defer consumeCancel()
	go sendAvailabilityRequestMessages(senderCtx, 2, time.Second)
	go kafka.Consume(consumeCtx, kafka.AvailabilityStatusRequestTopic, func(ctx context.Context, msg *kafka.GenericMessage) {
		asm, _ := kafka.NewAvailabilityStatusMessage(msg)
		require.EqualValues(t, "1", asm.SourceID)
		wg.Done()
	})

	msg, _ := kafka.AvailabilityStatusMessage{SourceID: "1"}.GenericMessage()
	EnqueueAvailabilityStatusRequest(&msg)
	EnqueueAvailabilityStatusRequest(&msg)
	time.Sleep(100 * time.Millisecond)
	senderCancel()
	wg.Wait()
}

func TestQueueCancelSend(t *testing.T) {
	ctx := context.Background()
	_ = kafka.InitializeStubBroker(16)

	wg := sync.WaitGroup{}
	wg.Add(1)
	consumeCtx, consumeCancel := context.WithCancel(ctx)
	senderCtx, senderCancel := context.WithCancel(ctx)
	defer consumeCancel()
	go sendAvailabilityRequestMessages(senderCtx, 2, 5*time.Second)
	go kafka.Consume(consumeCtx, kafka.AvailabilityStatusRequestTopic, func(ctx context.Context, msg *kafka.GenericMessage) {
		asm, _ := kafka.NewAvailabilityStatusMessage(msg)
		require.EqualValues(t, "1", asm.SourceID)
		wg.Done()
	})

	msg, _ := kafka.AvailabilityStatusMessage{SourceID: "1"}.GenericMessage()
	EnqueueAvailabilityStatusRequest(&msg)
	senderCancel()
	wg.Wait()
}
