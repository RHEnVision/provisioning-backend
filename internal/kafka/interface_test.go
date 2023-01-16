package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"

	"github.com/stretchr/testify/require"
)

func createMessage(topic, key, value string) *GenericMessage {
	return &GenericMessage{
		Topic: topic,
		Key:   []byte(key),
		Value: []byte(value),
	}
}

func TestSendAndConsume(t *testing.T) {
	ctx := context.Background()
	bus := NewStubBroker(16)
	m := createMessage("topic", "key", "value")

	wg := sync.WaitGroup{}
	wg.Add(3)
	cct, cancel := context.WithCancel(ctx)
	defer cancel()

	go bus.Consume(cct, "topic", time.Now(), func(ctx context.Context, msg *GenericMessage) {
		require.EqualValues(t, "key", msg.Key)
		require.EqualValues(t, "value", msg.Value)
		wg.Done()
	})

	_ = bus.Send(ctx, m, m, m)
	wg.Wait()
}

func TestMultipleTopics(t *testing.T) {
	ctx := context.Background()
	bus := NewStubBroker(16)

	wg := sync.WaitGroup{}
	wg.Add(2)
	cct, cancel := context.WithCancel(ctx)
	defer cancel()

	go bus.Consume(cct, "topic1", time.Now(), func(ctx context.Context, msg *GenericMessage) {
		require.EqualValues(t, "key1", msg.Key)
		require.EqualValues(t, "value", msg.Value)
		wg.Done()
	})

	go bus.Consume(cct, "topic2", time.Now(), func(ctx context.Context, msg *GenericMessage) {
		require.EqualValues(t, "key2", msg.Key)
		require.EqualValues(t, "value", msg.Value)
		wg.Done()
	})

	_ = bus.Send(ctx, createMessage("topic1", "key1", "value"))
	_ = bus.Send(ctx, createMessage("topic2", "key2", "value"))
	wg.Wait()
}
