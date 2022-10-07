package version

var (
	// Git SHA commit set via -ldflags
	BuildCommit string

	// Build date and time in UTC set via -ldflags
	BuildTime string
)

const (
	// PrometheusLabelName contains string that is included in all Prometheus metrics
	PrometheusLabelName = "provisioning"

	// APIPathName is the name used in main route API prefix
	APIPathName = PrometheusLabelName

	// APIPathName is the name used in main route API prefix
	APIPathVersion = "v1"

	// OpenTelemetryVersion is used for all OpenTelemetry tracing
	OpenTelemetryVersion = "1.0.0"
)
