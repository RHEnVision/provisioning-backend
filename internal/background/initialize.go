// Background jobs responsible for queueing batch operations (availability checks) or other
// operations (cleanups etc).
package background

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/rs/zerolog"
)

// Maximum batch size for each batch send, also an incoming buffered channel size to prevent
// incoming requests to overload the sender.
const availabilityStatusBatchSize = 16

// InitializeApi starts background goroutines for REST API processes.
// Use context cancellation to stop them.
func InitializeApi(ctx context.Context) {
	logger := zerolog.Ctx(ctx).With().Bool("background", true).Logger()
	ctx = logger.WithContext(ctx)

	// start availability request batch sender
	go sendAvailabilityRequestMessages(ctx, availabilityStatusBatchSize, 5*time.Second)
}

// InitializeWorker starts background goroutines for worker processes.
// Use context cancellation to stop them.
func InitializeWorker(ctx context.Context) {
	// no background goroutines at the moment
}

// InitializeStats starts background goroutines for the statuser process.
// Use context cancellation to stop it.
func InitializeStats(ctx context.Context) {
	logger := zerolog.Ctx(ctx).With().Bool("background", true).Logger()
	ctx = logger.WithContext(ctx)

	// start job queue telemetry
	go jobQueueMetricLoop(ctx, config.Stats.JobQueue, config.Hostname())

	// start database statistics
	go dbStatsLoop(ctx, config.Stats.ReservationsInterval)

	// cleanup old reservations
	if config.Reservation.CleanupEnabled {
		go dbCleanup(ctx, config.Reservation.CleanupInterval)
	}
}
