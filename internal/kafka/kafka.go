package kafka

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/random"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.opentelemetry.io/otel/trace"
)

type kafkaBroker struct {
	dialer    *kafka.Dialer
	transport *kafka.Transport
}

var _ Broker = &kafkaBroker{}

var (
	ErrDifferentTopic       = errors.New("messages in batch have different topics")
	ErrUnknownSASLMechanism = errors.New("unknown SASL mechanism")
)

func createSASLMechanism(saslMechanismName string, username string, password string) (sasl.Mechanism, error) {
	switch strings.ToLower(saslMechanismName) {
	case "plain", "none":
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
		return nil, fmt.Errorf("%w: %s", ErrUnknownSASLMechanism, saslMechanismName)
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
	var pool *x509.CertPool
	var tlsConfig *tls.Config
	var saslMechanism sasl.Mechanism

	logger := zerolog.Ctx(ctx)
	logger.Debug().Msgf("Setting up Kafka transport: %v", config.Kafka.Brokers)

	if config.Kafka.CACert != "" {
		logger.Debug().Str("cert", config.Kafka.CACert).Msg("Configuring TLS CA pool for Kafka")

		pemCerts := config.Kafka.CACert
		pool = x509.NewCertPool()
		if ok := pool.AppendCertsFromPEM([]byte(pemCerts)); !ok {
			logger.Warn().Msg("Could not add an CA cert to the pool")
		}
	}

	if config.Kafka.TlsEnabled && !config.InEphemeralClowder() {
		logger.Debug().Msg("Configuring Kafka for TLS")

		//nolint:gosec
		tlsConfig = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			RootCAs:            pool,
			InsecureSkipVerify: config.Kafka.TlsSkipVerify,
		}
	}

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

// kafka library has some noisy debug messages
var ignoredMsg *regexp.Regexp

func init() {
	var err error
	ignoredMsg, err = regexp.Compile("^no messages received from kafka within the allocated time.*")
	if err != nil {
		panic(err)
	}
}

func newContextLogger(ctx context.Context) func(msg string, a ...interface{}) {
	return func(msg string, a ...interface{}) {
		if ignoredMsg.MatchString(msg) {
			return
		}
		logger := zerolog.Ctx(ctx)
		logger.Debug().Bool("kafka", true).Msgf(msg, a...)
	}
}

func newContextErrLogger(ctx context.Context) func(msg string, a ...interface{}) {
	return func(msg string, a ...interface{}) {
		logger := zerolog.Ctx(ctx)
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
	logger := zerolog.Ctx(ctx)
	r := b.NewReader(ctx, topic)
	defer func() {
		if tempErr := r.Close(); tempErr != nil {
			logger.Warn().Err(tempErr).Msg("Unable to close the kafka reader")
		}
	}()

	err := r.SetOffsetAt(ctx, since)
	if err != nil {
		logger.Warn().Err(err).Msg("Unable to set initial offset")
	}

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil && errors.Is(err, io.EOF) {
			logger.Warn().Err(err).Msg("Kafka receiver has been closed")
			break
		} else if err != nil && errors.Is(err, context.Canceled) {
			logger.Debug().Msg("Kafka receiver has been cancelled")
			break
		} else if err != nil {
			logger.Warn().Err(err).Msg("Error when reading message")
		} else {
			logger.Trace().Bytes("payload", msg.Value).Msgf("Received message with key: %s, topic: %s, offset: %d, partition: %d",
				msg.Key, msg.Topic, msg.Offset, msg.Partition)

			// build new context - identity and trace id
			newLogger := logger.With().Str("msg_id", random.TraceID().String())
			newCtx, err := identity.WithIdentityFrom64(ctx, header("X-RH-Identity", msg.Headers))
			if err != nil {
				errLogger := newLogger.Logger()
				errLogger.Warn().Err(err).Msgf("Could not extract identity from context to Kafka message")
				newCtx = errLogger.WithContext(ctx)
			} else {
				id := identity.Identity(newCtx)

				traceId := trace.SpanFromContext(ctx).SpanContext().TraceID()
				if !traceId.IsValid() {
					traceId = random.TraceID()
				}
				newCtx = logging.WithTraceId(newCtx, traceId.String())

				newLogger = newLogger.
					Str("trace_id", traceId.String()).
					Str("account_number", id.Identity.AccountNumber).
					Str("org_id", id.Identity.OrgID)
				newCtx = newLogger.Logger().WithContext(newCtx)
			}

			handler(newCtx, NewMessageFromKafka(&msg))
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
// different topic than the first one, ErrDifferentTopic is returned.
func (b *kafkaBroker) Send(ctx context.Context, messages ...*GenericMessage) error {
	logger := zerolog.Ctx(ctx)

	if len(messages) == 0 {
		return nil
	}

	commonTopic := messages[0].Topic
	w := b.NewWriter(ctx)
	defer func() {
		if tempErr := w.Close(); tempErr != nil {
			logger.Warn().Err(tempErr).Msg("Unable to close the kafka writer")
		}
	}()

	logger.Trace().Str("topic", commonTopic).Msgf("Sending %d messages to Kafka", len(messages))

	kMessages := make([]kafka.Message, len(messages))
	for i, m := range messages {
		if m.Topic != commonTopic {
			return ErrDifferentTopic
		}
		kMessages[i] = m.KafkaMessage()
		if logger.Trace().Enabled() {
			dict := zerolog.Dict()
			for _, h := range kMessages[i].Headers {
				dict.Str(h.Key, string(h.Value))
			}
			logger.Trace().Str("topic", commonTopic).
				Bytes("body", kMessages[i].Value).
				Dict("headers", dict).
				Msg("Kafka message")
		}
	}

	err := w.WriteMessages(ctx, kMessages...)
	if err != nil {
		return fmt.Errorf("cannot send kafka messages(s): %w", err)
	}

	return nil
}
