package config

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/rs/zerolog"
)

// IsClowder returns true, when the app is running in clowder environment, whether it is
// production, stage or ephemeral.
func InClowder() bool {
	return clowder.IsClowderEnabled()
}

// InEphemeralClowder returns true, when the app is running in ephemeral clowder environment.
func InEphemeralClowder() bool {
	return clowder.IsClowderEnabled() && strings.HasPrefix(*clowder.LoadedConfig.Metadata.EnvName, "env-ephemeral")
}

func RedisHostAndPort() string {
	return fmt.Sprintf("%s:%d", Application.Cache.Redis.Host, Application.Cache.Redis.Port)
}

// InStageClowder returns true, when the app is running in stage clowder environment.
func InStageClowder() bool {
	return clowder.IsClowderEnabled() && strings.HasPrefix(*clowder.LoadedConfig.Metadata.EnvName, "env-stage")
}

// InProdClowder returns true, when the app is running in production clowder environment.
func InProdClowder() bool {
	return clowder.IsClowderEnabled() && strings.HasPrefix(*clowder.LoadedConfig.Metadata.EnvName, "env-prod")
}

// TopicName returns mapped topic from Clowder. When not running in Clowder mode, it returns the input topic name.
func TopicName(ctx context.Context, topic string) string {
	if t, ok := clowder.KafkaTopics[topic]; ok {
		return t.Name
	}
	if InClowder() {
		ctxval.Logger(ctx).Warn().Msgf("Tried to get TopicName for %s, but clowder doesn't know such topic", topic)
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
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
