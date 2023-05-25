package version

import "runtime/debug"

var (
	// Git SHA commit (first 4 characters)
	BuildCommit string

	// Build date and time
	BuildTime string

	// BuildGoVersion carries Go version the binary was built with
	BuildGoVersion string
)

const (
	// ApplicationName contains string that is used by other names
	ApplicationName = "provisioning"

	// PrometheusLabelName contains string that is included in all Prometheus metrics
	PrometheusLabelName = ApplicationName

	// APIPathName is the name used in main route API prefix
	APIPathName = ApplicationName

	// UnleashAppName contains string used to initialize Unleash (feature flags)
	UnleashAppName = ApplicationName

	// APIPathVersion is the name used in main route API prefix
	APIPathVersion = "v1"

	// OpenTelemetryVersion is used for all OpenTelemetry tracing
	OpenTelemetryVersion = "1.0.0"

	// KafkaClientID is used to identify Kafka client
	KafkaClientID = ApplicationName
)

func init() {
	bi, ok := debug.ReadBuildInfo()

	if !ok {
		BuildTime = "N/A"
		BuildCommit = "HEAD"
	}

	BuildGoVersion = bi.GoVersion

	for _, bs := range bi.Settings {
		switch bs.Key {
		case "vcs.revision":
			if len(bs.Value) > 4 {
				BuildCommit = bs.Value[0:4]
			}
		case "vcs.time":
			BuildTime = bs.Value
		}
	}
}
