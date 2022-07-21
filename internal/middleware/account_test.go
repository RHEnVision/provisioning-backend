package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/aws/smithy-go/ptr"
	"github.com/stretchr/testify/assert"
)

func TestAccountMiddleware(t *testing.T) {
	t.Run("existing account", func(t *testing.T) {
		ctx := identity.WithIdentity(t, context.Background())
		ctx = stubs.WithAccountDaoOne(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/v1/pubkeys", nil)
		if err != nil {
			assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))
		}

		rr := httptest.NewRecorder()

		isAccInNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var acc = ctxval.Account(r.Context())
			assert.NotNil(t, acc, "account was not set")
		})

		handler := middleware.AccountMiddleware(isAccInNext)
		handler.ServeHTTP(rr, req)
	})
	t.Run("create non-existing account", func(t *testing.T) {
		ctx := identity.WithCustomIdentity(t, context.Background(), "124", ptr.String("12"))
		ctx = stubs.WithAccountDaoOne(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/v1/pubkeys", nil)
		assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))

		rr := httptest.NewRecorder()

		isAccInNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var acc = ctxval.Account(r.Context())
			assert.NotNil(t, acc, "account was not set")
			assert.Equal(t, "124", acc.OrgID)
		})

		handler := middleware.AccountMiddleware(isAccInNext)
		handler.ServeHTTP(rr, req)

		count, err := stubs.AccountStubCount(ctx)
		assert.Nil(t, err, fmt.Sprintf("Error creating a new request: %v", err))
		assert.Equal(t, 2, count)
	})
}
