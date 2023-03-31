package telemetry

import (
	"context"
	"log"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// AppName is the application name which appears in traces.
const AppName = "provisioning-backend"

// TracePrefix must be used as prefix for TraceNames in packages.
const TracePrefix = AppName + "/"

type Telemetry struct {
	tracerProvider *trace.TracerProvider
	propagator     *propagation.TextMapPropagator
}

func Middleware(routes chi.Routes) func(next http.Handler) http.Handler {
	return otelchi.Middleware(AppName, otelchi.WithChiRoutes(routes), otelchi.WithRequestMethodInSpanName(true))
}

func Initialize(rootLogger *zerolog.Logger) *Telemetry {
	if config.Telemetry.Type == "" {
		return &Telemetry{}
	}
	logger := rootLogger.With().Bool("otel", true).Logger()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("provisioning"),
		semconv.ServiceVersionKey.String(version.OpenTelemetryVersion),
	)

	var tpOptions []trace.TracerProviderOption
	if config.Telemetry.Type == "sentry" {
		// production use case: sentry exporting
		sentryLogger := logger.With().Bool("sentry", true).Logger()
		sentryLogger.Debug().Msg("Initializing otel sentry integration")
		logWriter := logging.DebugWriter{Logger: &sentryLogger}
		sentry.Logger = log.New(&logWriter, "", 0)
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.Sentry.Dsn,
			EnableTracing:    true,
			TracesSampleRate: config.Sentry.SampleRate,
			Debug:            config.Sentry.Debug,
			DebugWriter:      &logWriter,
			HTTPTransport: &LoggingTransport{
				roundTripper: http.DefaultTransport,
				logger:       &sentryLogger,
			},
		})
		if err != nil {
			panic(err)
		}

		tpOptions = append(tpOptions, trace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()))
	} else if config.Telemetry.Type == "jaeger" {
		// production use case: full exporting, batching
		exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Telemetry.Jaeger.Endpoint)))
		if err != nil {
			panic(err)
		}
		tpOptions = append(tpOptions, trace.WithResource(res))
		tpOptions = append(tpOptions, trace.WithBatcher(exporter))
	} else if config.Telemetry.Type == "logger" {
		// development use case: logger exporting, synchronous
		exporter := NewZerologExporter(&logger)
		tpOptions = append(tpOptions, trace.WithResource(res))
		tpOptions = append(tpOptions, trace.WithSyncer(exporter))
	} else {
		// No tracing configured - do nothing
		exporter := NewNoopExporter()
		tpOptions = append(tpOptions, trace.WithResource(res))
		tpOptions = append(tpOptions, trace.WithSyncer(exporter))
	}

	tp := trace.NewTracerProvider(tpOptions...)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, sentryotel.NewSentryPropagator())
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagator)
	return &Telemetry{tracerProvider: tp, propagator: &propagator}
}

func (t *Telemetry) Close(ctx context.Context) {
	if t.tracerProvider == nil {
		return
	}
	_ = t.tracerProvider.Shutdown(context.Background())
}
