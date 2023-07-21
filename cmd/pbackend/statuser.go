package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/notifications"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/random"

	// Clients
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/azure"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/image_builder"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/sources"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/kafka"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"

	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const ChannelBuffer = 32

type SourceInfo struct {
	MessageContext      context.Context // Carries logger and identity
	Authentication      clients.Authentication
	SourceApplicationID string
}

var (
	chAws        = make(chan SourceInfo, ChannelBuffer)
	chAzure      = make(chan SourceInfo, ChannelBuffer)
	chGcp        = make(chan SourceInfo, ChannelBuffer)
	chSend       = make(chan kafka.SourceResult, ChannelBuffer)
	receiverWG   = sync.WaitGroup{}
	processingWG = sync.WaitGroup{}
	senderWG     = sync.WaitGroup{}
)

func init() {
	random.SeedGlobal()
}

func processMessage(msgCtx context.Context, message *kafka.GenericMessage) {
	logger := zerolog.Ctx(msgCtx)

	// Get source id
	asm, err := kafka.NewAvailabilityStatusMessage(message)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get availability status message")
		return
	}

	// Set source id as logging field
	sourceId := asm.SourceID
	logger = ptr.To(logger.With().Str("source_id", sourceId).Logger())
	ctx := logger.WithContext(msgCtx)
	logger.Trace().Msgf("Sources availability check for %s", sourceId)

	// Get sources client
	sourcesClient, err := clients.GetSourcesClient(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get sources client")
		return
	}

	// Fetch authentication from Sources
	authentication, err := sourcesClient.GetAuthentication(ctx, sourceId)
	if err != nil {
		metrics.IncTotalInvalidAvailabilityCheckReqs()
		if errors.Is(err, clients.NotFoundErr) {
			logger.Warn().Err(err).Msg("Not found error from sources")
			return
		}
		logger.Warn().Err(err).Msg("Could not get authentication")
		return
	}

	s := SourceInfo{
		MessageContext:      ctx,
		Authentication:      *authentication,
		SourceApplicationID: authentication.SourceApplictionID,
	}

	switch authentication.ProviderType {
	case models.ProviderTypeAWS:
		chAws <- s
	case models.ProviderTypeAzure:
		chAzure <- s
	case models.ProviderTypeGCP:
		chGcp <- s
	case models.ProviderTypeNoop:
	case models.ProviderTypeUnknown:
		logger.Warn().Err(err).Msg("Authentication provider type is unknown")
	}
}

func checkSourceAvailabilityAzure(cancelCtx context.Context) {
	defer processingWG.Done()

	for s := range chAzure {
		ctx := s.MessageContext
		logger := zerolog.Ctx(ctx)

		if random.Float32() > config.Azure.AvailabilityRate {
			logger.Trace().Msgf("Skipping Azure source availability status %s", s.SourceApplicationID)
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeAzure.String(), "skipped", nil)
			continue
		}

		logger.Trace().Msgf("Checking Azure source availability status %s", s.SourceApplicationID)
		metrics.ObserveAvailabilityCheckReqsDuration(models.ProviderTypeAzure.String(), func() error {
			var err error
			sr := kafka.SourceResult{
				MessageContext: ctx,
				ResourceID:     s.SourceApplicationID,
				ResourceType:   "Application",
			}
			sr.Status = kafka.StatusAvailable
			chSend <- sr
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeAzure.String(), sr.Status.String(), nil)

			return fmt.Errorf("error during check: %w", err)
		})

		if cancelCtx.Err() != nil {
			break
		}

		time.Sleep(config.Azure.AvailabilityDelay)
	}
}

func checkSourceAvailabilityAWS(cancelCtx context.Context) {
	defer processingWG.Done()

	for s := range chAws {
		ctx := s.MessageContext
		logger := zerolog.Ctx(ctx)

		if random.Float32() > config.AWS.AvailabilityRate {
			logger.Trace().Msgf("Skipping AWS source availability status %s", s.SourceApplicationID)
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeAWS.String(), "skipped", nil)
			continue
		}

		logger.Trace().Msgf("Checking AWS source availability status %s", s.SourceApplicationID)
		metrics.ObserveAvailabilityCheckReqsDuration(models.ProviderTypeAWS.String(), func() error {
			var err error
			var permissions []string
			sr := kafka.SourceResult{
				MessageContext: ctx,
				ResourceID:     s.SourceApplicationID,
				ResourceType:   "Application",
			}
			ec2Client, err := clients.GetEC2Client(ctx, &s.Authentication, "")
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Err(err).Msg("Could not get aws assumed client")
			} else {
				sr.Status = kafka.StatusAvailable
				permissions, err = ec2Client.CheckPermission(ctx, &s.Authentication)
				if err != nil {
					sr.Status = kafka.StatusUnavailable
					sr.Err = err
					sr.MissingPermissions = permissions
					if logger.Info().Enabled() {
						arr := zerolog.Arr()
						for _, p := range permissions {
							arr.Str(p)
						}
						logger.Info().Err(err).
							Array("missing_aws_permissions", arr).
							Str("source_id", sr.ResourceID).Msg("Missing AWS permissions")
					}
				}
			}
			chSend <- sr
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeAWS.String(), sr.Status.String(), err)
			return fmt.Errorf("error during check: %w", err)
		})

		if cancelCtx.Err() != nil {
			break
		}

		time.Sleep(config.AWS.AvailabilityDelay)
	}
}

func checkSourceAvailabilityGCP(cancelCtx context.Context) {
	defer processingWG.Done()

	for s := range chGcp {
		ctx := s.MessageContext
		logger := zerolog.Ctx(ctx)

		if random.Float32() > config.GCP.AvailabilityRate {
			logger.Trace().Msgf("Skipping GCP source availability status %s", s.SourceApplicationID)
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeGCP.String(), "skipped", nil)
			continue
		}

		logger.Trace().Msgf("Checking GCP source availability status %s", s.SourceApplicationID)
		metrics.ObserveAvailabilityCheckReqsDuration(models.ProviderTypeGCP.String(), func() error {
			var err error
			sr := kafka.SourceResult{
				MessageContext: s.MessageContext,
				ResourceID:     s.SourceApplicationID,
				ResourceType:   "Application",
			}
			gcpClient, err := clients.GetGCPClient(ctx, &s.Authentication)
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Err(err).Msg("Could not get gcp client")
				chSend <- sr
			}
			_, err = gcpClient.ListAllRegions(ctx)
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Err(err).Msg("Could not list gcp regions")
				chSend <- sr
			} else {
				sr.Status = kafka.StatusAvailable
				chSend <- sr
			}
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeGCP.String(), sr.Status.String(), err)

			return fmt.Errorf("error during check: %w", err)
		})

		if cancelCtx.Err() != nil {
			break
		}

		time.Sleep(config.GCP.AvailabilityDelay)
	}
}

func sendResults(cancelCtx context.Context, batchSize int, tickDuration time.Duration) {
	messages := make([]*kafka.GenericMessage, 0, batchSize)
	ticker := time.NewTicker(tickDuration)
	defer senderWG.Done()

	for {
		select {

		case sr := <-chSend:
			ctx := sr.MessageContext
			logger := zerolog.Ctx(ctx)
			msg, err := sr.GenericMessage(ctx)
			if err != nil {
				logger.Warn().Err(err).Msg("Could not generate generic message")
				continue
			}
			messages = append(messages, &msg)
			length := len(messages)

			if length >= batchSize {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (full buffer)", length)
				err := kafka.Send(ctx, messages...)
				if err != nil {
					logger.Warn().Err(err).Msg("Could not send source availability status messages (full buffer)")
				}
				messages = messages[:0]
			}
		case <-ticker.C:
			logger := zerolog.Ctx(cancelCtx)
			length := len(messages)
			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (tick)", length)
				err := kafka.Send(cancelCtx, messages...)
				if err != nil {
					logger.Warn().Err(err).Msg("Could not send source availability status messages (tick)")
				}
				messages = messages[:0]
			}
		case <-cancelCtx.Done():
			logger := zerolog.Ctx(cancelCtx)
			ticker.Stop()
			length := len(messages)

			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (cancel)", length)
				err := kafka.Send(cancelCtx, messages...)
				if err != nil {
					logger.Warn().Err(err).Msg("Could not send source availability status messages (cancel)")
				}
			}

			return
		}
	}
}

func statuser() {
	ctx := context.Background()
	config.Initialize("config/api.env", "config/statuser.env")

	// initialize cloudwatch using the AWS clients
	logger, closeFunc := logging.InitializeLogger()
	defer closeFunc()
	logging.DumpConfigForDevelopment()

	// initialize telemetry
	tel := telemetry.Initialize(&log.Logger)
	defer tel.Close(ctx)

	// initialize platform kafka and notifications
	if config.Kafka.Enabled {
		err := kafka.InitializeKafkaBroker(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("Unable to initialize the platform kafka")
		}

		if config.Application.Notifications.Enabled {
			notifications.Initialize(ctx)
		}
	}

	// metrics
	logger.Info().Msgf("Starting new instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle(config.Prometheus.Path, promhttp.Handler())
	metricsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Prometheus.Port),
		Handler: metricsRouter,
	}

	signalNotify := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		if shutdownErr := metricsServer.Shutdown(context.Background()); shutdownErr != nil {
			logger.Warn().Err(shutdownErr).Msg("Metrics service shutdown error")
		}
		close(signalNotify)
	}()

	go func() {
		if listenErr := metricsServer.ListenAndServe(); listenErr != nil {
			var errInUse syscall.Errno
			if errors.As(listenErr, &errInUse) && errInUse == syscall.EADDRINUSE {
				logger.Warn().Err(listenErr).Msg("Not starting metrics service, port already in use")
			} else if !errors.Is(listenErr, http.ErrServerClosed) {
				logger.Warn().Err(listenErr).Msg("Metrics service listen error")
			}
		}
	}()

	// start the consumer
	receiverWG.Add(1)
	cancelCtx, consumerCancelFunc := context.WithCancel(ctx)
	consumerNotify := make(chan struct{})
	go func() {
		defer receiverWG.Done()
		kafka.Consume(cancelCtx, kafka.AvailabilityStatusRequestTopic, time.Now(), processMessage)
		close(consumerNotify)
	}()

	metrics.RegisterStatuserMetrics()

	// initialize the database
	logger.Debug().Msg("Initializing database connection")
	err := db.Initialize(ctx, "public")
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}
	defer db.Close()

	// start processing goroutines
	processingWG.Add(3)

	go checkSourceAvailabilityAWS(cancelCtx)
	go checkSourceAvailabilityGCP(cancelCtx)
	go checkSourceAvailabilityAzure(cancelCtx)

	senderWG.Add(1)
	// processing can be slowed down by configuration: send messages in maximum
	// batch size of 8 messages or 4 seconds (what comes first)
	go sendResults(cancelCtx, 8, 4*time.Second)

	logger.Info().Msg("Statuser process started")
	select {
	case <-signalNotify:
		logger.Info().Msg("Exiting due to signal")
	case <-consumerNotify:
		logger.Warn().Msg("Exiting due to closed consumer")
	}

	// stop kafka receiver (can take up to 10 seconds) and wait until it returns
	consumerCancelFunc()
	receiverWG.Wait()

	// close all processors and wait until it exits the range loop
	close(chAws)
	close(chAzure)
	close(chGcp)
	processingWG.Wait()

	// close the sending channel and wait until it exits the range loop
	close(chSend)
	senderWG.Wait()

	logger.Info().Msg("Consumer shutdown initiated")
	consumerCancelFunc()
	logger.Info().Msg("Shutdown finished, exiting")
}
