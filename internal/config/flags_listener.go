package config

import (
	"github.com/Unleash/unleash-client-go/v3"
	"github.com/rs/zerolog"
)

type logListener struct {
	logger *zerolog.Logger
}

// OnError prints out errors.
func (l logListener) OnError(err error) {
	l.logger.Error().Err(err).Msg("Unleash error")
}

// OnWarning prints out warning.
func (l logListener) OnWarning(err error) {
	l.logger.Warn().Err(err).Msg("Unleash warning")
}

// OnReady prints to the console when the repository is ready.
func (l logListener) OnReady() {
	l.logger.Info().Msg("Unleash ready")
}

// OnCount prints to the console when the feature is queried.
func (l logListener) OnCount(name string, enabled bool) {
	l.logger.Trace().Msgf("Unleash query '%s': %t", name, enabled)
}

// OnSent prints to the console when the server has uploaded metrics.
func (l logListener) OnSent(payload unleash.MetricsData) {
	l.logger.Trace().Msgf("Unleash metrics data: %+v", payload)
}

// OnRegistered prints to the console when the client has registered.
func (l logListener) OnRegistered(payload unleash.ClientData) {
	l.logger.Trace().Msgf("Registered unleash client: %+v", payload)
}
