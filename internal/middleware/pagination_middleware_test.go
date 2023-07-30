package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaginationMiddleware(t *testing.T) {
	t.Run("with limit and offset", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "/api/data?offset=10&limit=20", nil)
		require.NoError(t, err, "failed to create request")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limit := page.Limit(r.Context()).Int()
			offset := page.Offset(r.Context()).Int()
			assert.Equal(t, 10, offset)
			assert.Equal(t, 20, limit)
		})

		paginationHandler := Pagination(handler)
		paginationHandler.ServeHTTP(rr, req)
	})
	t.Run("without limit and offset", func(t *testing.T) {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "/api/data?offset=&limit=", nil)
		require.NoError(t, err, "failed to create request")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limit := page.Limit(r.Context()).Int()
			offset := page.Offset(r.Context()).Int()
			assert.Equal(t, 0, offset)
			assert.Equal(t, 100, limit)
		})

		paginationHandler := Pagination(handler)
		paginationHandler.ServeHTTP(rr, req)
	})
}
