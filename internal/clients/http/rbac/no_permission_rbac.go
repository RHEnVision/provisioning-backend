package rbac

import (
	"context"
	"errors"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type noPermRbac struct{}

var noPermRbacSingleton = noPermRbac{}

// GetPrincipalAccess returns ACL that has no permissions. Used when Rbac connection could not be established.
func (r noPermRbac) GetPrincipalAccess(ctx context.Context) (clients.RbacAcl, error) {
	return clients.NoPermissionsRbacAcl, nil
}

var ErrNoPermissionRbac = errors.New("RBAC client could not be established")

// Ready returns an error that RBAC could not be established
func (r noPermRbac) Ready(ctx context.Context) error {
	return ErrNoPermissionRbac
}
