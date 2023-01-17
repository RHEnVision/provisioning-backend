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
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
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

	SourceID string

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

func processMessage(ctx context.Context, message *kafka.GenericMessage) {
	logger := ctxval.Logger(ctx)

	// Get source id
	asm, err := kafka.NewAvailabilityStatusMessage(message)
	if err != nil {
		logger.Warn().Msgf("could not get availability status message %s", err)
		return
	}

	sourceId := asm.SourceID

	// Get sources client
	sourcesClient, err := clients.GetSourcesClient(ctx)
	if err != nil {
		logger.Warn().Msgf("Could not get sources client %s", err)
		return
	}

	// Fetch authentication from Sources
	authentication, err := sourcesClient.GetAuthentication(ctx, sourceId)
	if err != nil {
		if errors.Is(err, clients.NotFoundErr) {
			logger.Warn().Msgf("Not found error: %s", err)
			return
		}
		logger.Warn().Msgf("Could not get authentication: %s", err)
		return
	}

	s := SourceInfo{
		Authentication: *authentication,
		SourceID:       sourceId,
		Identity:       ctxval.Identity(ctx),
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
		logger.Warn().Msg("Authentication provider type is unknown")
	}
}

func checkSourceAvailabilityAzure(ctx context.Context) {
	defer processingWG.Done()

	for s := range chAzure {
		metrics.ObserveAvailablilityCheckReqsDuration(models.ProviderTypeAzure, func() error {
			var err error
			sr := kafka.SourceResult{
				SourceID:     s.SourceID,
				Identity:     s.Identity,
				ResourceType: "Source",
			}
			// TODO: check if source is avavliable - WIP
			sr.Status = kafka.StatusAvaliable
			chSend <- sr
			metrics.IncTotalAvailabilityCheckReqs(models.ProviderTypeAzure, sr.Status, nil)

			return fmt.Errorf("error during check %w:", err)
		})
	}
}

func checkSourceAvailabilityAWS(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	defer processingWG.Done()

	for s := range chAws {
		metrics.ObserveAvailablilityCheckReqsDuration(models.ProviderTypeAWS, func() error {
			var err error
			sr := kafka.SourceResult{
				SourceID:     s.SourceID,
				Identity:     s.Identity,
				ResourceType: "Source",
			}
			_, err = clients.GetEC2Client(ctx, &s.Authentication, "")
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Msgf("Could not get aws assumed client %s", err)
				chSend <- sr
			} else {
				sr.Status = kafka.StatusAvaliable
				chSend <- sr
			}
			metrics.IncTotalAvailabilityCheckReqs(models.ProviderTypeAWS, sr.Status, err)
			return fmt.Errorf("error during check %w:", err)
		})
	}
}

func checkSourceAvailabilityGCP(ctx context.Context) {
	logger := ctxval.Logger(ctx)
	defer processingWG.Done()

	for s := range chGcp {
		metrics.ObserveAvailablilityCheckReqsDuration(models.ProviderTypeGCP, func() error {
			var err error
			sr := kafka.SourceResult{
				SourceID:     s.SourceID,
				Identity:     s.Identity,
				ResourceType: "Source",
			}
			gcpClient, err := clients.GetGCPClient(ctx, &s.Authentication)
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Msgf("Could not get gcp client %s", err)
				chSend <- sr
			}
			_, err = gcpClient.ListAllRegions(ctx)
			if err != nil {
				sr.Status = kafka.StatusUnavailable
				sr.Err = err
				logger.Warn().Msgf("Could not list gcp regions %s", err)
				chSend <- sr
			} else {
				sr.Status = kafka.StatusAvaliable
				chSend <- sr
			}
			metrics.IncTotalAvailabilityCheckReqs(models.ProviderTypeGCP, sr.Status, err)

			return fmt.Errorf("error during check %w:", err)
		})
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
				logger.Warn().Msgf("Could not generate generic message %s", err)
				continue
			}
			messages = append(messages, &msg)
			length := len(messages)

			if length >= batchSize {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (full buffer)", length)
				err := kafka.Send(ctx, messages...)
				if err != nil {
					logger.Warn().Msgf("Could not send source availability status messages (full buffer) %s", err)
				}
				messages = messages[:0]
			}
		case <-ticker.C:
			length := len(messages)
			if length > 0 {
				logger.Trace().Int("messages", length).Msgf("Sending %d source availability status messages (tick)", length)
				err := kafka.Send(ctx, messages...)
				if err != nil {
					logger.Warn().Msgf("Could not send source availability status messages (tick) %s", err)
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
					logger.Warn().Msgf("Could not send source availability status messages (cancel) %s", err)
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
