package jobs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// copyContext returns a new context with key values copied. Used when the original
// context has expired but there is still some work to be done.
func copyContext(ctx context.Context) context.Context {
	nCtx := context.Background()
	nCtx = log.Logger.WithContext(nCtx)
	nCtx = logging.WithTraceId(nCtx, logging.TraceId(ctx))
	nCtx = logging.WithEdgeRequestId(nCtx, logging.EdgeRequestId(ctx))
	nCtx = identity.WithAccountId(nCtx, identity.AccountId(ctx))
	return nCtx
}

func reservationContextLogger(origCtx context.Context, reservationID int64) (context.Context, *zerolog.Logger) {
	ctx := logging.WithReservationId(origCtx, reservationID)

	logger := zerolog.Ctx(ctx)
	logger = ptr.To(logger.With().
		Int64("reservation_id", reservationID).
		Logger())
	ctx = logger.WithContext(ctx)
	return ctx, logger
}
