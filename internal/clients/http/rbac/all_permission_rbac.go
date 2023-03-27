package rbac

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type allPermRbac struct{}

var allPermRbacSingleton = allPermRbac{}

// GetPrincipalAccess returns ACL that has all permissions. Used when Rbac is disabled by the configuration.
func (r allPermRbac) GetPrincipalAccess(ctx context.Context) (clients.RbacAcl, error) {
	return clients.AllPermissionsRbacAcl, nil
}

// Ready returns no error
func (r allPermRbac) Ready(ctx context.Context) error {
	return nil
}
