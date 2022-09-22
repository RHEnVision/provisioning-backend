package http

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const TraceName = "github.com/EnVision/provisioning/internal/clients/http"

// OpenTelemetryDoer gathers OpenTelemetry tracing for every request. For a given
// name, the HTTP method of the request is concatenated as a span name:
//
// imagebuilder-get
//
// URL is not part of span name as this would create too many spans. When error occurs,
// it is properly propagated into the span. It is recommended to create additional
// spans for each function in the client to see how much time is spent doing the actual
// request and how much time a function need to process the data.
type OpenTelemetryDoer struct {
	name string
	ctx  context.Context
	doer HttpRequestDoer
}

func NewOpenTelemetryDoer(ctx context.Context, doer HttpRequestDoer, name string) *OpenTelemetryDoer {
	return &OpenTelemetryDoer{
		name: name,
		ctx:  ctx,
		doer: doer,
	}
}

func (doer *OpenTelemetryDoer) Do(req *http.Request) (*http.Response, error) {
	spanName := doer.name + "-" + req.Method
	ctx, span := otel.Tracer(TraceName).Start(req.Context(), spanName)
	defer span.End()

	// delegate the request
	resp, err := doer.doer.Do(req.WithContext(ctx))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, NewDoerErr(err)
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return resp, nil
}
