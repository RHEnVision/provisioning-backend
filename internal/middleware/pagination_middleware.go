package middleware

import (
	"net/http"

	"github.com/RHEnVision/provisioning-backend/internal/page"
)

// Pagination middleware is used to extract the offset and the limit
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("offset")
		limit := r.URL.Query().Get("limit")

		newCtx := page.WithOffset(r.Context(), offset)
		newCtx = page.WithLimit(newCtx, limit)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}
