package telemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	otrace "go.opentelemetry.io/otel/trace"
)

const AppName = "provisioning-backend"

type Telemetry struct {
	tracerProvider *trace.TracerProvider
	propagator     propagation.TextMapPropagator
}

func Middleware(routes chi.Routes) func(next http.Handler) http.Handler {
	return otelchi.Middleware(AppName, otelchi.WithChiRoutes(routes), otelchi.WithRequestMethodInSpanName(true))
}

func Initialize(ctx context.Context, rootLogger *zerolog.Logger) *Telemetry {
	if !config.Telemetry.Enabled {
		return &Telemetry{}
	}
	logger := rootLogger.With().Bool("otel", true).Logger()

	var exporterOption trace.TracerProviderOption
	if config.Telemetry.OTLP.Enabled {
		logger.Debug().Msg("Jaeger telemetry enabled")
		endpoint := fmt.Sprintf("%s:%d", config.Telemetry.OTLP.Hostname, config.Telemetry.OTLP.Port)
		// production use case: full exporting, batching, compression
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(endpoint),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		}
		if config.Telemetry.OTLP.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		exporter, err := otlptracehttp.New(ctx, opts...)
		if err != nil {
			panic(err)
		}
		exporterOption = trace.WithBatcher(exporter)
	} else if config.Telemetry.Logger.Enabled {
		logger.Debug().Msg("Logger telemetry enabled")
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
	otel.SetErrorHandler(&ZerologOpenTelemetryErrorHandler{logger: &logger})
	return &Telemetry{tracerProvider: tp, propagator: propagator}
}

func (t *Telemetry) Close(_ context.Context) {
	if t.tracerProvider == nil {
		return
	}
	_ = t.tracerProvider.Shutdown(context.Background())
}

// StartSpan wraps starting a tracing span.
// This makes sure we start spans only when Telemetry is enabled.
//
// It also wraps fetching the Tracer, which is a bit confusing,
// but its name is just for distinguishing different Tracers,
// we have not identified need for multiple tracers in our code yet.
// We use empty name, it should be safe according to
// https://pkg.go.dev/go.opentelemetry.io/otel/trace@v1.16.0#TracerProvider
func StartSpan(ctx context.Context, spanName string, opts ...otrace.SpanStartOption) (context.Context, otrace.Span) {
	if config.Telemetry.Enabled {
		return otel.Tracer("").Start(ctx, spanName, opts...)
	} else {
		// return empty invalid span
		return ctx, otrace.SpanFromContext(context.Background())
	}
}
