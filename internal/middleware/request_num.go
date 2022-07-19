package middleware

import (
	"net/http"
	"sync/atomic"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
)

var reqNum uint64

func RequestNum(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddUint64(&reqNum, 1)
		ctx := ctxval.WithRequestNum(r.Context(), n)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
