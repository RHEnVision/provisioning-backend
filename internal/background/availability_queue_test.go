package background

import (
	"context"
	"runtime"
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

	// enqueue message to be sent first
	msg, _ := kafka.AvailabilityStatusMessage{SourceID: "1"}.GenericMessage()
	EnqueueAvailabilityStatusRequest(&msg)

	// start sending messages
	senderCtx, senderCancel := context.WithCancel(ctx)
	go sendAvailabilityRequestMessages(senderCtx, 2, 5*time.Second)

	// allow the other goroutine to put the message into the buffer
	runtime.Gosched()

	// set the receiving message function up
	wg := sync.WaitGroup{}
	wg.Add(1)
	consumeCtx, consumeCancel := context.WithCancel(ctx)
	defer consumeCancel()
	go kafka.Consume(consumeCtx, kafka.AvailabilityStatusRequestTopic, func(ctx context.Context, msg *kafka.GenericMessage) {
		asm, _ := kafka.NewAvailabilityStatusMessage(msg)
		require.EqualValues(t, "1", asm.SourceID)
		wg.Done()
	})

	// at this point hopefully the message was already buffered but just in case :-)
	time.Sleep(20 * time.Millisecond)

	// cancel the sender before the 5 seconds timeout (so cancel branch is triggered)
	senderCancel()

	// wait until the message is consumed
	wg.Wait()
}
