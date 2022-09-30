package telemetry

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
)

type loggerExporter struct {
	logger *zerolog.Logger
}

func NewZerologExporter(logger *zerolog.Logger) trace.SpanExporter {
	return &loggerExporter{logger: logger}
}

func (e *loggerExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	for _, span := range spans {
		traceId := span.SpanContext().TraceID().String()
		spanId := span.SpanContext().SpanID().String()
		duration := span.EndTime().Sub(span.StartTime())
		name := span.Name()
		parentId := span.Parent().SpanID()
		statusCode := span.Status().Code
		statusMsg := span.Status().Description
		t := e.logger.Trace().
			Str("trace_id", traceId).
			Str("span_id", spanId).
			Dur("duration", duration)

		if parentId.IsValid() {
			t = t.Str("span_id_parent", parentId.String())
		}

		if statusCode != codes.Unset {
			t = t.Str("status_code", statusCode.String())
		}

		if len(statusMsg) > 0 {
			t = t.Str("status_msg", statusMsg)
		}

		for _, kv := range span.Attributes() {
			key := string(kv.Key)
			val := kv.Value.AsString()

			if val == "" {
				continue
			}

			if key == "msg" || key == "message" {
				key = "otel_msg"
			}

			t = t.Str(key, val)
		}
		t.Msgf("otel: %s", name)
	}
	return nil
}

func (e *loggerExporter) Shutdown(ctx context.Context) error {
	return nil
}
