package ctxval

import "context"

// Copy returns a new context with key values copied. Used when the original
// context has expired but there is still some work to be done.
func Copy(ctx context.Context) context.Context {
	nCtx := context.Background()
	nCtx = WithLogger(nCtx, Logger(ctx))
	nCtx = WithTraceId(nCtx, TraceId(ctx))
	nCtx = WithEdgeRequestId(nCtx, EdgeRequestId(ctx))
	nCtx = WithAccountId(nCtx, AccountId(ctx))
	return nCtx
}
