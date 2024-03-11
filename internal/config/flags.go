package config

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/version"
	"github.com/Unleash/unleash-client-go/v4"
	"github.com/rs/zerolog"
)

const unleashProjectName = "default"

// FeatureEnabled returns a user-specific or global feature flag. When such flag is not found,
// true is returned.
func FeatureEnabled(ctx context.Context, name string) bool {
	// When unleash is not configured, all features are enabled by default.
	if !config.Unleash.Enabled {
		return true
	}

	uctx := UnleashContext(ctx)
	return unleash.IsEnabled(name, unleash.WithContext(uctx), unleash.WithFallback(true))
}

// LaunchEnabled: launch button in UI initiates the application workflow
func LaunchEnabled(ctx context.Context) bool {
	return FeatureEnabled(ctx, fmt.Sprintf("%s.launch", config.Unleash.Prefix))
}

func unleashLogger(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx).With().Bool("unleash", true).Logger()
	return &logger
}

// InitializeFeatureFlags configures unleash client and starts poller routine. Callers
// must close poller by calling StopFeatureFlags function.
func InitializeFeatureFlags(ctx context.Context) error {
	if !config.Unleash.Enabled {
		return nil
	}

	listener := logListener{
		logger: unleashLogger(ctx),
	}
	err := unleash.Initialize(
		unleash.WithListener(&listener),
		unleash.WithUrl(config.Unleash.URL),
		unleash.WithProjectName(unleashProjectName),
		unleash.WithAppName(version.UnleashAppName),
		unleash.WithEnvironment(config.Unleash.Environment),
		unleash.WithCustomHeaders(http.Header{"Authorization": {config.Unleash.Token}}),
	)
	if err != nil {
		return fmt.Errorf("unleash error: %w", err)
	}

	return nil
}

// StopFeatureFlags stops the unleash feature flag poller
func StopFeatureFlags(ctx context.Context) {
	if !config.Unleash.Enabled {
		return
	}

	err := unleash.Close()
	if err != nil {
		unleashLogger(ctx).Warn().Err(err).Msg("Unable to close unleash poller")
	}
}
