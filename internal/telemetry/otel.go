package telemetry

import (
	"context"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/version"
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

const AppName = "provisioning-backend"

const TracePrefix = AppName + "/"

type Telemetry struct {
	tracerProvider *trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

func Middleware(routes chi.Routes) func(next http.Handler) http.Handler {
	return otelchi.Middleware(AppName, otelchi.WithChiRoutes(routes), otelchi.WithRequestMethodInSpanName(true))
}

func Initialize(rootLogger *zerolog.Logger) *Telemetry {
	if !config.Telemetry.Enabled {
		return &Telemetry{}
	}
	logger := rootLogger.With().Bool("otel", true).Logger()

	var exporterOption trace.TracerProviderOption
	if config.Telemetry.Jaeger.Enabled {
		// production use case: full exporting, batching
		exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.Telemetry.Jaeger.Endpoint)))
		if err != nil {
			panic(err)
		}
		exporterOption = trace.WithBatcher(exporter)
	} else if config.Telemetry.Logger.Enabled {
		// development use case: logger exporting, synchronous
		exporter := NewZerologExporter(&logger)
		exporterOption = trace.WithSyncer(exporter)
	} else {
		// No tracing configured - do nothing
		exporter := NewNoopExporter()
		exporterOption = trace.WithSyncer(exporter)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("provisioning"),
		semconv.ServiceVersionKey.String(version.OpenTelemetryVersion),
	)

	tp := trace.NewTracerProvider(
		exporterOption,
		trace.WithResource(res),
	)

	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagator)
	return &Telemetry{tracerProvider: tp, propagator: propagator}
}

func (t *Telemetry) Close(_ context.Context) {
	if t.tracerProvider == nil {
		return
	}
	_ = t.tracerProvider.Shutdown(context.Background())
}
