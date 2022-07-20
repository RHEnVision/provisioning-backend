package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"

	"github.com/rs/xid"
)

var RequestIDHeader = "X-Request-Id"

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = xid.New().String()
		}
		ctx = ctxval.SetRequestId(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
