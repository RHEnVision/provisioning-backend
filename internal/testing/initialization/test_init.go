// The package initializes logging and all stubs for use in tests.
package initialization

import (
	// Initialize logging (must be kept the first)
	_ "github.com/RHEnVision/provisioning-backend/internal/logging/testing"

	// HTTP client stub implementations
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"

	// DAO stub implementation
	_ "github.com/RHEnVision/provisioning-backend/internal/dao/stubs"

	// Job queue stub
	_ "github.com/RHEnVision/provisioning-backend/internal/jobs/queue/stub"

	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger.Debug().Msg("initialized logging, HTTP client(s), DAO and job queue stubs")
}
