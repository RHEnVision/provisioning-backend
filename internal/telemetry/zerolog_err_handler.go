package telemetry

import (
	"github.com/rs/zerolog"
)

type ZerologOpenTelemetryErrorHandler struct {
	logger *zerolog.Logger
}

func (eh ZerologOpenTelemetryErrorHandler) Handle(err error) {
	eh.logger.Warn().Err(err).Msgf("OTLP error: %s", err.Error())
}
