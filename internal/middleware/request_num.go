package middleware

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"net/http"
	"sync/atomic"
)

var reqNum uint64

func RequestNum(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		n := atomic.AddUint64(&reqNum, 1)
		ctx = context.WithValue(ctx, ctxval.RequestNumCtxKey, n)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
