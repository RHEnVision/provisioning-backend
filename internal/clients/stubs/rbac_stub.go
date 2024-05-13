package stubs

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

type rbacClient struct{}

func getRbacClient(ctx context.Context) clients.Rbac {
	return rbacClient{}
}

func (r rbacClient) GetPrincipalAccess(ctx context.Context) (clients.RbacAcl, error) {
	return clients.AllPermissionsRbacAcl, nil
}

func (r rbacClient) Ready(ctx context.Context) error {
	return nil
}
