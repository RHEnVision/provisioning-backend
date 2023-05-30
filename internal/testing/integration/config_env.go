package integration

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/zerolog/log"
)

func InitConfigEnvironment(ctx context.Context, envPath string) context.Context {
	config.Initialize("config/test.env", envPath)
	logging.InitializeStdout()
	return log.Logger.WithContext(ctx)
}
