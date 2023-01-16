package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type kafkaBroker struct {
	dialer    *kafka.Dialer
	transport *kafka.Transport
}

var _ Broker = &kafkaBroker{}

var (
	DifferentTopicErr       = errors.New("messages in batch have different topics")
	UnknownSaslMechanismErr = errors.New("unknown SASL mechanism")
)

func createSASLMechanism(saslMechanismName string, username string, password string) (sasl.Mechanism, error) {
	switch strings.ToLower(saslMechanismName) {
	case "plain":
		return plain.Mechanism{
			Username: username,
			Password: password,
		}, nil
	case "scram-sha-512":
		mechanism, err := scram.Mechanism(scram.SHA512, username, password)
		if err != nil {
			return nil, fmt.Errorf("unable to create scram-sha-512 mechanism: %w", err)
		}

		return mechanism, nil
	case "scram-sha-256":
		mechanism, err := scram.Mechanism(scram.SHA256, username, password)
		if err != nil {
			return nil, fmt.Errorf("unable to create scram-sha-256 mechanism: %w", err)
		}

		return mechanism, nil
	default:
		return nil, fmt.Errorf("%w: %s", UnknownSaslMechanismErr, saslMechanismName)
	}
}

func InitializeKafkaBroker(ctx context.Context) error {
	var err error
	broker, err = NewKafkaBroker(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize kafka: %w", err)
	}

	InitializeTopicRequests(ctx)

	return nil
}

func NewKafkaBroker(ctx context.Context) (Broker, error) {
	var tlsConfig *tls.Config
	var saslMechanism sasl.Mechanism

	logger := ctxval.Logger(ctx)
	logger.Debug().Msgf("Setting up Kafka transport: %v CA:%v SASL:%v", config.Kafka.Brokers,
		config.Kafka.CACert != "", config.Kafka.SASL.SaslMechanism != "" && config.Kafka.SASL.SaslMechanism != "none")

	// configure TLS when CA certificate was provided
	if config.Kafka.CACert != "" {
		logger.Debug().Str("cert", config.Kafka.CACert).Msg("Adding CA certificates to the pool")

		pemCerts := config.Kafka.CACert
		pool := x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(pemCerts)); !ok {
			logger.Warn().Msg("Could not add an CA cert to the pool")
		}

		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
			RootCAs:    pool,
		}
	}

	// configure SASL if mechanism was provided
	if config.Kafka.SASL.SaslMechanism != "" {
		var err error
		saslMechanism, err = createSASLMechanism(config.Kafka.SASL.SaslMechanism, config.Kafka.SASL.Username, config.Kafka.SASL.Password)
		if err != nil {
			return nil, fmt.Errorf("kafka SASL error: %w", err)
		}
	}

	dialer := &kafka.Dialer{
		ClientID:      version.KafkaClientID,
		Timeout:       10 * time.Second,
		SASLMechanism: saslMechanism,
		TLS:           tlsConfig,
	}

	transport := &kafka.Transport{
		Dial:     dialer.DialFunc,
		ClientID: version.KafkaClientID,
		TLS:      tlsConfig,
		SASL:     saslMechanism,
	}

	return &kafkaBroker{
		dialer:    dialer,
		transport: transport,
	}, nil
}

func newContextLogger(ctx context.Context) func(msg string, a ...interface{}) {
	return func(msg string, a ...interface{}) {
		logger := ctxval.Logger(ctx)
		logger.Debug().Bool("kafka", true).Msgf(msg, a...)
	}
}

func newContextErrLogger(ctx context.Context) func(msg string, a ...interface{}) {
	return func(msg string, a ...interface{}) {
		logger := ctxval.Logger(ctx)
		logger.Warn().Bool("kafka", true).Msgf(msg, a...)
	}
}

// NewReader creates a reader. Use Close() function to close the reader.
func (b *kafkaBroker) NewReader(ctx context.Context, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Kafka.Brokers,
		Dialer:      b.dialer,
		Topic:       topic,
		StartOffset: kafka.LastOffset,
		Logger:      kafka.LoggerFunc(newContextLogger(ctx)),
		ErrorLogger: kafka.LoggerFunc(newContextErrLogger(ctx)),
	})
}

// NewWriter creates synchronous writer created from the pool. It does not have associated any topic with it,
// therefore topic must be set on the message-level. Make sure to close it with Close() function.
func (b *kafkaBroker) NewWriter(ctx context.Context) *kafka.Writer {
	return &kafka.Writer{
		Addr:        kafka.TCP(config.Kafka.Brokers...),
		Transport:   b.transport,
		Logger:      kafka.LoggerFunc(newContextLogger(ctx)),
		ErrorLogger: kafka.LoggerFunc(newContextErrLogger(ctx)),
	}
}

// Consume reads messages in batches up to 1 MB with up to 10 seconds delay. It blocks, therefore
// it should be called from a separate goroutine. Use context cancellation to stop the loop.
func (b *kafkaBroker) Consume(ctx context.Context, topic string, since time.Time, handler func(ctx context.Context, message *GenericMessage)) {
	logger := ctxval.Logger(ctx)
	r := b.NewReader(ctx, topic)
	defer r.Close()

	err := r.SetOffsetAt(ctx, since)
	if err != nil {
		logger.Warn().Err(err).Msg("Unable to set initial offset")
	}

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil && errors.Is(err, io.EOF) {
			break
		} else if err != nil && errors.Is(err, context.Canceled) {
			break
		} else if err != nil {
			logger.Warn().Err(err).Msgf("Error when reading message: %s", err.Error())
		} else {
			logger.Trace().Bytes("payload", msg.Value).Msgf("Received message with key: %s, topic: %s, offset: %d, partition: %d",
				msg.Key, msg.Topic, msg.Offset, msg.Partition)
			ctx, err = ctxval.WithIdentityFrom64(ctx, header("x-rh-identity", msg.Headers))
			if err != nil {
				logger.Trace().Msgf("Could not extract identity from context to Kafka message: %s", err)
			}
			handler(ctx, NewMessageFromKafka(&msg))
		}
	}
}

func header(name string, headers []protocol.Header) string {
	for _, h := range headers {
		if strings.EqualFold(h.Key, name) {
			return string(h.Value)
		}
	}
	return ""
}

// Send one or more generic messages with the same topic. If there is a message with
// different topic than the first one, DifferentTopicErr is returned.
func (b *kafkaBroker) Send(ctx context.Context, messages ...*GenericMessage) error {
	logger := ctxval.Logger(ctx)

	if len(messages) == 0 {
		return nil
	}

	commonTopic := messages[0].Topic
	w := b.NewWriter(ctx)
	defer w.Close()

	logger.Trace().Str("topic", commonTopic).Msgf("Sending %d messages to Kafka", len(messages))

	kMessages := make([]kafka.Message, len(messages))
	for i, m := range messages {
		if m.Topic != commonTopic {
			return DifferentTopicErr
		}
		kMessages[i] = m.KafkaMessage()
	}

	err := w.WriteMessages(ctx, kMessages...)
	if err != nil {
		return fmt.Errorf("cannot send kafka messages(s): %w", err)
	}

	return nil
}
