package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/middleware"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	tidentity "github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	_ "github.com/RHEnVision/provisioning-backend/internal/testing/initialization"
	"github.com/stretchr/testify/assert"
)

func TestAccountMiddleware(t *testing.T) {
	t.Run("existing account", func(t *testing.T) {
		ctx := stubs.WithAccountDaoOne(context.Background())
		ctx = tidentity.WithTenant(t, ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/test", nil)
		if err != nil {
			t.Errorf("Error creating a test request: %v", err)
		}

		rr := httptest.NewRecorder()

		isAccInNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acc := identity.AccountId(r.Context())
			assert.NotNil(t, acc, "account id was not set")
		})

		handler := middleware.AccountMiddleware(isAccInNext)
		handler.ServeHTTP(rr, req)
	})
	t.Run("create non-existing account", func(t *testing.T) {
		ctx := tidentity.WithCustomIdentity(t, context.Background(), "124", ptr.To("12"))
		ctx = stubs.WithAccountDaoOne(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/test", nil)
		if err != nil {
			t.Errorf("Error creating a test request: %v", err)
		}

		rr := httptest.NewRecorder()

		isAccInNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accId := identity.AccountId(r.Context())
			assert.NotNil(t, accId, "account id was not set")
			accDao := dao.GetAccountDao(r.Context())
			acc, innerErr := accDao.GetById(r.Context(), accId)
			if innerErr != nil {
				t.Errorf("could not fetch account by id: %v", err)
			}
			assert.Equal(t, "124", acc.OrgID)
		})

		handler := middleware.AccountMiddleware(isAccInNext)
		handler.ServeHTTP(rr, req)

		count := stubs.AccountStubCount(ctx)
		assert.Equal(t, 2, count)
	})
	t.Run("existing null account", func(t *testing.T) {
		ctx := stubs.WithAccountDaoNull(context.Background())
		ctx = tidentity.WithTenant(t, ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/v1/pubkeys", nil)
		if err != nil {
			require.NoError(t, err, fmt.Sprintf("Error creating a new request: %v", err))
		}

		rr := httptest.NewRecorder()

		isAccInNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			acc := identity.AccountId(r.Context())
			assert.NotNil(t, acc, "account id was not set")
		})

		handler := middleware.AccountMiddleware(isAccInNext)
		handler.ServeHTTP(rr, req)
	})
	t.Run("create non-existing null account", func(t *testing.T) {
		ctx := tidentity.WithCustomIdentity(t, context.Background(), "124", ptr.To("12"))
		ctx = stubs.WithAccountDaoNull(ctx)

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/provisioning/v1/pubkeys", nil)
		require.NoError(t, err, fmt.Sprintf("Error creating a new request: %v", err))

		rr := httptest.NewRecorder()

		isAccInNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accId := identity.AccountId(r.Context())
			assert.NotNil(t, accId, "account id was not set")
			accDao := dao.GetAccountDao(r.Context())
			acc, innerErr := accDao.GetById(r.Context(), accId)
			if innerErr != nil {
				t.Errorf("could not fetch account by id: %v", err)
			}
			assert.Equal(t, "124", acc.OrgID)
		})

		handler := middleware.AccountMiddleware(isAccInNext)
		handler.ServeHTTP(rr, req)

		count := stubs.AccountStubCount(ctx)
		assert.Equal(t, 2, count)
	})
}
