package config

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/rs/zerolog"
)

// IsClowder returns true, when the app is running in clowder environment, whether it is
// production, stage or ephemeral.
func InClowder() bool {
	return clowder.IsClowderEnabled()
}

// Environment returns clowder environment (stage, prod or ephemeral) or "dev" in case the app
// is not running in Clowder.
func Environment() string {
	if InClowder() && clowder.LoadedConfig != nil && clowder.LoadedConfig.Metadata != nil && clowder.LoadedConfig.Metadata.EnvName != nil {
		return *clowder.LoadedConfig.Metadata.EnvName
	}

	return "dev"
}

// EnvironmentPrefix wraps an identifier (e.g. id, unique id) in the following way.
// For production environment, it returns "prefix-identifier". For any other environment,
// it returns "prefix-identifier-env". Examples:
//
// * reservation-14572
// * reservation-9531-stage
// * reservation-13-ephemeral
// * reservation-1-dev
func EnvironmentPrefix(prefix, identifier string) string {
	env := Environment()

	if strings.HasPrefix(env, "prod") || env == "" {
		return prefix + "-" + identifier
	}

	return prefix + "-" + identifier + "-" + env
}

// InEphemeralClowder returns true, when the app is running in ephemeral clowder environment.
func InEphemeralClowder() bool {
	return InClowder() && strings.Contains(*clowder.LoadedConfig.Metadata.EnvName, "ephemeral")
}

func RedisHostAndPort() string {
	return fmt.Sprintf("%s:%d", Application.Cache.Redis.Host, Application.Cache.Redis.Port)
}

// InStageClowder returns true, when the app is running in stage clowder environment.
func InStageClowder() bool {
	return InClowder() && strings.Contains(*clowder.LoadedConfig.Metadata.EnvName, "stage")
}

// InProdClowder returns true, when the app is running in production clowder environment.
func InProdClowder() bool {
	return InClowder() && strings.Contains(*clowder.LoadedConfig.Metadata.EnvName, "prod")
}

// TopicName returns mapped topic from Clowder. When not running in Clowder mode, it returns the input topic name.
func TopicName(ctx context.Context, topic string) string {
	if t, ok := clowder.KafkaTopics[topic]; ok {
		return t.Name
	}
	if InClowder() {
		zerolog.Ctx(ctx).Warn().Msgf("Tried to get TopicName for %s, but clowder doesn't know such topic", topic)
	}

	return topic
}

func StringToURL(urlStr string) *url.URL {
	if urlStr == "" {
		return nil
	}
	urlProxy, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}
	return urlProxy
}

// DumpClowder writes safely some information from clowder config.
func DumpClowder(logger zerolog.Logger) {
	if !InClowder() {
		return
	}
	cfg := clowder.LoadedConfig
	brokers := make([]string, len(cfg.Kafka.Brokers))
	for i, b := range cfg.Kafka.Brokers {
		brokers[i] = b.Hostname
	}

	logger.Info().Msgf("Clowder environment: %s", *cfg.Metadata.EnvName)
	logger.Info().Msgf("Clowder database hostname: %s", cfg.Database.Hostname)
	logger.Info().Msgf("Clowder kafka brokers: %s", strings.Join(brokers, ","))
	logger.Info().Msgf("Clowder logging type: %s, region: %s, group: %s",
		cfg.Logging.Type,
		cfg.Logging.Cloudwatch.Region,
		cfg.Logging.Cloudwatch.LogGroup)
}

// DumpConfig writes configuration to a logger. It removes all secrets, however, it is never
// recommended to call this function in production environments.
func DumpConfig(logger zerolog.Logger) {
	replacement := "****"
	configCopy := config
	configCopy.Database.Password = replacement
	configCopy.Cloudwatch.Key = replacement
	configCopy.Cloudwatch.Secret = replacement
	configCopy.Cloudwatch.Session = replacement
	configCopy.AWS.Key = replacement
	configCopy.AWS.Secret = replacement
	configCopy.AWS.Session = replacement
	configCopy.RestEndpoints.Sources.Password = replacement
	configCopy.RestEndpoints.ImageBuilder.Password = replacement
	configCopy.Azure.ClientID = replacement
	configCopy.Azure.ClientSecret = replacement
	configCopy.GCP.JSON = replacement
	configCopy.Unleash.Token = replacement
	configCopy.Kafka.SASL.Username = replacement
	configCopy.Kafka.SASL.Password = replacement
	// We want to know if the DSN was empty
	if configCopy.Sentry.Dsn != "" {
		configCopy.Sentry.Dsn = replacement
	}
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
