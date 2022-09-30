# OpenTelemetry

While logging is done through Zerolog with ability to export to AWS CloudWatch, the project supports tracing via OpenTelemetry and handover of W3C Trace Context.

# Configuration

To enable OpenTelemetry, use the following configuration:

```yaml
telemetry:
  enabled: true
```

By default, OpenTelemetry data is not exported because full OpenTelemetry Jaeger tracing support is planned for the platform. As a temporary solution, tho, it is possible to export tracing data into the application logger (Zerolog/CloudWatch):

```yaml
telemetry:
    enabled: true
    loggerExporter: true

logging:
    level: -1
```

Keep in mind that logging level must be set to "trace" (-1) as all the tracing data is reported at this level. All the OpenTelemetry log messages start with `otel:` and also have "otel" boolean flag set to true. Span attributes are exported as regular Zerolog fields. Example output:

```
11:19AM TRC otel: /api/provisioning/v1/ready/{SRV} duration=6.689334 hostname=mone.home.lan http.flavor=1.1 http.host=localhost:8000 http.method=GET http.route=/api/provisioning/v1/ready/{SRV} http.scheme=http http.server_name=provisioning-backend http.target=/api/provisioning/v1/ready/sources http.user_agent="Apache-HttpClient/4.5.13 (Java/17.0.3)" net.host.name=localhost net.peer.ip=127.0.0.1 net.transport=ip_tcp otel=true span_id=0a557df38618c977 trace_id=6c87886d623eb0a9c5e3ec03bb64edfb
11:19AM TRC otel: HTTP GET duration=2418.531833 hostname=mone.home.lan http.flavor=1.1 http.host=console.stage.redhat.com http.method=GET http.scheme=http http.url=https://console.stage.redhat.com/api/image-builder... otel=true span_id=8fb56982443c1ae0 span_id_parent=a9c1040b9686ba5d trace_id=a58e792bb7754c7132e9812374cf0243
11:19AM TRC otel: Ready duration=2419.063708 hostname=mone.home.lan otel=true span_id=72638bd6a72a23ed span_id_parent=0b9434ba2311ac5f trace_id=a58e792bb7754c7132e9812374cf0243
11:19AM INF Completed GET request /api/provisioning/v1/ready/ib in 2.491s ms with 200 bytes_in=0 bytes_out=0 hostname=mone.home.lan latency_ms=2490.723792 method=GET remote_ip=127.0.0.1:50180 status=200 trace_id=a58e792bb7754c7132e9812374cf0243 url=/api/provisioning/v1/ready/ib
11:19AM TRC otel: /api/provisioning/v1/ready/{SRV} duration=2490.834875 hostname=mone.home.lan http.flavor=1.1 http.host=localhost:8000 http.method=GET http.route=/api/provisioning/v1/ready/{SRV} http.scheme=http http.server_name=provisioning-backend http.target=/api/provisioning/v1/ready/ib http.user_agent="Apache-HttpClient/4.5.13 (Java/17.0.3)" net.host.name=localhost net.peer.ip=127.0.0.1 net.transport=ip_tcp otel=true span_id=0b9434ba2311ac5f trace_id=a58e792bb7754c7132e9812374cf0243
```

## Features

Tracing ID is parsed from the W3C Trace Context header or generated when missing for each incoming request. The Trace ID is generated even if tracing feature is turned off because this field is used for correlation of log messages for each request on the application level.

Spans are created for each Chi route with the route being the name of the span (e.g. `/api/provisioning/v1/ready/{SRV}`).

Spans are created for each HTTP client call being made via `telemetry.HTTPClient`. Name is set to "HTTP" and followed by HTTP method.

Spans are created for all SQL operations made through the `pgx` SQL driver.

Spans are created for custom instrumentation points. An example:

```go
func Function() {
    ctx, span := otel.Tracer(TraceName).Start(ctx, "Function")
    defer span.End()
    // ...
}
```

