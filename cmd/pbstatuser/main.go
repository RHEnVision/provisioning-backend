package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/random"
	"github.com/getsentry/sentry-go"
	"github.com/redhatinsights/platform-go-middlewares/identity"

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
	"github.com/rs/zerolog/log"
)

const ChannelBuffer = 32

type SourceInfo struct {
	Authentication clients.Authentication

	SourceApplicationID string

	Identity identity.XRHID
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

func processMessage(origCtx context.Context, message *kafka.GenericMessage) {
	logger := ctxval.Logger(origCtx)

	// Get source id
	asm, err := kafka.NewAvailabilityStatusMessage(message)
	if err != nil {
		logger.Warn().Err(err).Msgf("Could not get availability status message %s", err)
		return
	}
	logger.Trace().Msgf("Received a message from sources to be processed with source id %s", asm.SourceID)

	sourceId := asm.SourceID

	// Set source id as logging field
	logger = ptr.To(logger.With().Str("source_id", sourceId).Logger())
	ctx := ctxval.WithLogger(origCtx, logger)

	// Get sources client
	sourcesClient, err := clients.GetSourcesClient(ctx)
	if err != nil {
		logger.Warn().Err(err).Msgf("Could not get sources client %s", err)
		return
	}

	// Fetch authentication from Sources
	authentication, err := sourcesClient.GetAuthentication(ctx, sourceId)
	if err != nil {
		metrics.IncTotalInvalidAvailabilityCheckReqs()
		if errors.Is(err, clients.NotFoundErr) {
			logger.Warn().Err(err).Msgf("Not found error: %s", err)
			return
		}
		logger.Warn().Err(err).Msgf("Could not get authentication: %s", err)
		return
	}

	s := SourceInfo{
		Authentication:      *authentication,
		SourceApplicationID: authentication.SourceApplictionID,
		Identity:            ctxval.Identity(ctx),
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

func checkSourceAvailabilityAzure(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	defer processingWG.Done()

	for s := range chAzure {
		logger.Trace().Msgf("Checking Azure source availability status %s", s.SourceApplicationID)
		var err error
		metrics.ObserveAvailabilityCheckReqsDuration(models.ProviderTypeAzure, func() error {
			sr := kafka.SourceResult{
				ResourceID:   s.SourceApplicationID,
				Identity:     s.Identity,
				ResourceType: "Application",
			}
			// TODO: check if source is available - WIP
			sr.Status = kafka.StatusAvaliable
			chSend <- sr
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeAzure, sr.Status, nil)

			return err
		})
		if err != nil && config.Sentry.Enabled {
			sentry.CaptureException(err)
		}
	}
}

func checkSourceAvailabilityAWS(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	defer processingWG.Done()

	for s := range chAws {
		logger.Trace().Msgf("Checking AWS source availability status %s", s.SourceApplicationID)
		var err error
		metrics.ObserveAvailabilityCheckReqsDuration(models.ProviderTypeAWS, func() error {
			sr := kafka.SourceResult{
				ResourceID:   s.SourceApplicationID,
				Identity:     s.Identity,
				ResourceType: "Application",
			}
			_, err = clients.GetEC2Client(ctx, &s.Authentication, "")
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Err(err).Msgf("Could not get aws assumed client %s", err)
				chSend <- sr
				err = fmt.Errorf("error during check: %w", err)
			} else {
				sr.Status = kafka.StatusAvaliable
				chSend <- sr
			}
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeAWS, sr.Status, err)
			return err
		})
		if err != nil && config.Sentry.Enabled {
			sentry.CaptureException(err)
		}
	}
}

func checkSourceAvailabilityGCP(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	defer processingWG.Done()

	for s := range chGcp {
		logger.Trace().Msgf("Checking GCP source availability status %s", s.SourceApplicationID)
		var err error
		metrics.ObserveAvailabilityCheckReqsDuration(models.ProviderTypeGCP, func() error {
			sr := kafka.SourceResult{
				ResourceID:   s.SourceApplicationID,
				Identity:     s.Identity,
				ResourceType: "Application",
			}
			gcpClient, gcpErr := clients.GetGCPClient(ctx, &s.Authentication)
			if gcpErr != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = gcpErr
				logger.Warn().Err(gcpErr).Msgf("Could not get gcp client %s", gcpErr)
				chSend <- sr
				err = fmt.Errorf("error during check: %w", gcpErr)
			}
			_, gcpErr = gcpClient.ListAllRegions(ctx)
			if gcpErr != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = gcpErr
				logger.Warn().Err(gcpErr).Msgf("Could not list gcp regions %s", gcpErr)
				chSend <- sr
				err = fmt.Errorf("error during check: %w", gcpErr)
			} else {
				sr.Status = kafka.StatusAvaliable
				chSend <- sr
			}
			metrics.IncTotalSentAvailabilityCheckReqs(models.ProviderTypeGCP, sr.Status, err)

			return err
		})
		if err != nil && config.Sentry.Enabled {
			sentry.CaptureException(err)
		}
	}
}

func sendResults(ctx context.Context, batchSize int, tickDuration time.Duration) {
	logger := ctxval.Logger(ctx)
	messages := make([]*kafka.GenericMessage, 0, batchSize)
	ticker := time.NewTicker(tickDuration)
	defer senderWG.Done()

	for {
		select {

		case sr := <-chSend:
			ctx = ctxval.WithIdentity(ctx, sr.Identity)
			msg, err := sr.GenericMessage(ctx)
			if err != nil {
				logger.Warn().Err(err).Msgf("Could not generate generic message %s", err)
				if config.Sentry.Enabled {
					sentry.CaptureException(err)
				}
				continue
			}
			messages = append(messages, &msg)
			length := len(messages)

			if length >= batchSize {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (full buffer)", length)
				err = kafka.Send(ctx, messages...)
				if err != nil {
					logger.Warn().Err(err).Msgf("Could not send source availability status messages (full buffer) %s", err)
					if config.Sentry.Enabled {
						sentry.CaptureException(err)
					}
				}
				messages = messages[:0]
			}
		case <-ticker.C:
			length := len(messages)
			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (tick)", length)
				err := kafka.Send(ctx, messages...)
				if err != nil {
					logger.Warn().Err(err).Msgf("Could not send source availability status messages (tick) %s", err)
					if config.Sentry.Enabled {
						sentry.CaptureException(err)
					}
				}
				messages = messages[:0]
			}
		case <-ctx.Done():
			ticker.Stop()
			length := len(messages)

			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (cancel)", length)
				err := kafka.Send(ctx, messages...)
				if err != nil {
					logger.Warn().Err(err).Msgf("Could not send source availability status messages (cancel) %s", err)
					if config.Sentry.Enabled {
						sentry.CaptureException(err)
					}
				}
			}

			return
		}
	}
}

func main() {
	ctx := context.Background()
	config.Initialize("config/api.env", "config/statuser.env")

	// initialize stdout logging and AWS clients first
	logging.InitializeStdout()

	// initialize cloudwatch using the AWS clients
	logger, clsFunc, err := logging.InitializeCloudwatch(log.Logger)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing cloudwatch")
	}
	defer clsFunc()
	log.Logger = logger
	logging.DumpConfigForDevelopment()

	// setup structured logging
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-hostname"
	}
	logger = logger.With().
		Timestamp().
		Str("hostname", hostname).
		Bool("statuser", true).
		Logger()

	// initialize Sentry error logging
	if config.Sentry.Enabled {
		sentry.Init(sentry.ClientOptions{
			Dsn: config.Sentry.Dsn,
		})
		sentry.ConfigureScope(func(scope *sentry.Scope) {
			scope.SetTag("stream", path.Base(os.Args[0]))
		})
		defer sentry.Recover()
	}

	// initialize telemetry
	tel := telemetry.Initialize(&log.Logger)
	defer tel.Close(ctx)

	// initialize platform kafka
	logger.Info().Msg("Initializing platform kafka")
	err = kafka.InitializeKafkaBroker(ctx)
	if err != nil {
		logger.Fatal().Err(err).Msg("Unable to initialize the platform kafka")
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
		if err := metricsServer.Shutdown(context.Background()); err != nil {
			logger.Fatal().Err(err).Msg("Metrics service shutdown error")
		}
		close(signalNotify)
	}()

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			var errInUse syscall.Errno
			if errors.As(err, &errInUse) && errInUse == syscall.EADDRINUSE {
				logger.Warn().Err(err).Msg("Not starting metrics service, port already in use")
			} else if !errors.Is(err, http.ErrServerClosed) {
				logger.Warn().Err(err).Msg("Metrics service listen error")
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

	// start processing goroutines
	processingWG.Add(3)

	go checkSourceAvailabilityAWS(cancelCtx)
	go checkSourceAvailabilityGCP(cancelCtx)
	go checkSourceAvailabilityAzure(cancelCtx)

	senderWG.Add(1)
	go sendResults(cancelCtx, 1024, 5*time.Second)

	logger.Info().Msg("Worker started")
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
