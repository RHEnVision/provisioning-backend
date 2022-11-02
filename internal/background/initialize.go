// Background jobs responsible for queueing batch operations (availability checks) or other
// operations (cleanups etc).
package background

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

// Maximum batch size for each batch send, also an incoming buffered channel size to prevent
// incoming requests to overload the sender.
const availabilityStatusBatchSize = 1024

// Initialize starts background goroutines. Use context cancellation to stop them.
func Initialize(ctx context.Context) {
	logger := ctxval.Logger(ctx).With().Bool("background", true).Logger()
	ctx = ctxval.WithLogger(ctx, &logger)

	// start availability request batch sender
	go sendAvailabilityRequestMessages(ctx, availabilityStatusBatchSize, 5*time.Second)
}
