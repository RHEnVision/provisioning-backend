package testing

import (
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = logging.InitializeStdout()
}
