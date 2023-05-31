package jobs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
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
