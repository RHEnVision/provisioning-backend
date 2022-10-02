package testing

import (
	"github.com/RHEnVision/provisioning-backend/internal/logging"
)

func init() {
	logging.InitializeStdout()
}
