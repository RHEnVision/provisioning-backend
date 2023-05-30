package azure

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/telemetry"
	"github.com/rs/zerolog"
)

const TraceName = telemetry.TracePrefix + "internal/clients/http/azure"

func logger(ctx context.Context) zerolog.Logger {
	return zerolog.Ctx(ctx).With().Str("client", "azure").Logger()
}
