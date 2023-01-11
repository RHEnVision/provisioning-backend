package kafka

import (
	"context"
	"sync"
)

// In-memory broker
type stubBroker struct {
	data       map[string]chan *GenericMessage
	m          sync.Mutex
	bufferSize int
}

var _ Broker = &stubBroker{}

func InitializeStubBroker(bufferSize int) error {
	broker = NewStubBroker(bufferSize)

	return nil
}

func NewStubBroker(bufferSize int) Broker {
	return &stubBroker{
		data:       make(map[string]chan *GenericMessage),
		bufferSize: bufferSize,
	}
}

func (s *stubBroker) find(topic string) chan *GenericMessage {
	s.m.Lock()
	defer s.m.Unlock()

	if ch, ok := s.data[topic]; ok {
		return ch
	} else {
		ch := make(chan *GenericMessage, s.bufferSize)
		s.data[topic] = ch
		return ch
	}
}

func (s *stubBroker) Consume(ctx context.Context, topic string, group string, handler func(ctx context.Context, message *GenericMessage)) {
	ch := s.find(topic)

	for {
		select {
		case msg := <-ch:
			handler(ctx, msg)
		case <-ctx.Done():
			return
		}
	}
}

func (s *stubBroker) Send(_ context.Context, messages ...*GenericMessage) error {
	for _, m := range messages {
		ch := s.find(m.Topic)
		ch <- m
	}

	return nil
}
