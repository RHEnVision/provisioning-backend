package config

import (
	"net/url"

	clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/rs/zerolog"
)

// IsClowder returns true, when the app is running in clowder environment, whether it is
// production, stage or ephemeral.
func InClowder() bool {
	return clowder.IsClowderEnabled()
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

// DumpConfig writes configuration to a logger. It removes all secrets, however, it is never
// recommended to call this function in production environments.
func DumpConfig(logger zerolog.Logger) {
	if InClowder() {
		logger.Warn().Msg("Dumping configuration in production mode!")
	}
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
	logger.Info().Msgf("Configuration: %+v", configCopy)
}
