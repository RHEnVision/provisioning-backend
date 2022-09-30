package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/sdk/trace"
)

type exporter struct{}

func NewNoopExporter() trace.SpanExporter {
	return &exporter{}
}

func (_ *exporter) ExportSpans(_ context.Context, _ []trace.ReadOnlySpan) error {
	return nil
}

func (_ *exporter) Shutdown(_ context.Context) error {
	return nil
}
