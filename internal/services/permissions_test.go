package services

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/rbac"
	"github.com/stretchr/testify/require"
)

func TestCheckPermissionAndRenderJoin(t *testing.T) {
	ctx := context.Background()
	ctx = rbac.WithAcl(ctx, clients.AccessList{
		clients.NewAccess("provisioning:r1.r2:perm"),
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	require.NoError(t, err, "failed to create request")

	err = CheckPermissionAndRender(w, req, "perm", "r1", "r2")
	require.NoError(t, err)
}

func TestCheckPermissionAndRenderEmpty(t *testing.T) {
	ctx := context.Background()
	ctx = rbac.WithAcl(ctx, clients.AccessList{
		clients.NewAccess("provisioning:r1:perm"),
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	require.NoError(t, err, "failed to create request")

	err = CheckPermissionAndRender(w, req, "perm", "r1", "")
	require.NoError(t, err)
}

func TestCheckPermissionAndRenderOne(t *testing.T) {
	ctx := context.Background()
	ctx = rbac.WithAcl(ctx, clients.AccessList{
		clients.NewAccess("provisioning:r1:perm"),
	})

	w := httptest.NewRecorder()
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	require.NoError(t, err, "failed to create request")

	err = CheckPermissionAndRender(w, req, "perm", "r1")
	require.NoError(t, err)
}
