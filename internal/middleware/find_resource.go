package middleware

import (
	"context"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"net/http"
)

func FindResourceCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxval.ResourceCtxKey, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
