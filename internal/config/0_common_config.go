package config

import "github.com/RHEnVision/provisioning-backend/internal/config/parser"

// Common initialization must be called first.
func init() {
	// App
	parser.Viper.SetDefault("app.name", "provisioning")
	parser.Viper.SetDefault("app.port", 8000)
	parser.Viper.SetDefault("app.version", "v1")

	// Database
	parser.Viper.SetDefault("database.host", "localhost")
	parser.Viper.SetDefault("database.port", 5432)
	parser.Viper.SetDefault("database.name", "")
	parser.Viper.SetDefault("database.user", "")
	parser.Viper.SetDefault("database.password", "")
	parser.Viper.SetDefault("database.loglevel", 1)

	// Logging
	parser.Viper.SetDefault("logging.level", 1)

	// Cloudwatch
	parser.Viper.SetDefault("cloudwatch.enabled", false)
	parser.Viper.SetDefault("cloudwatch.region", "")
	parser.Viper.SetDefault("cloudwatch.key", "")
	parser.Viper.SetDefault("cloudwatch.secret", "")
	parser.Viper.SetDefault("cloudwatch.session", "")
	parser.Viper.SetDefault("cloudwatch.group", "")
	parser.Viper.SetDefault("cloudwatch.stream", "")

	// Prometheus
	parser.Viper.SetDefault("prometheus.port", 9000)
	parser.Viper.SetDefault("prometheus.path", "/metrics")

	// FeatureFlags
	parser.Viper.SetDefault("featureFlags.environment", "production")

	// RestEndpoints
	parser.Viper.SetDefault("restEndpoints.imageBuilder.url", "")
	parser.Viper.SetDefault("restEndpoints.sources.url", "")
	parser.Viper.SetDefault("restEndpoints.traceData", false)

	// Worker
	parser.Viper.SetDefault("worker.queue", "memory")
	parser.Viper.SetDefault("worker.concurrency", 8)
	parser.Viper.SetDefault("worker.heartbeatSec", 30)
	parser.Viper.SetDefault("worker.maxBeats", 10)

	// AWS
	parser.Viper.SetDefault("aws.key", "")
	parser.Viper.SetDefault("aws.secret", "")
	parser.Viper.SetDefault("aws.session", "")
	parser.Viper.SetDefault("aws.instancePrefix", "")
}
