package middleware

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"net/http"

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
		ctx = context.WithValue(ctx, ctxval.RequestIdCtxKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
