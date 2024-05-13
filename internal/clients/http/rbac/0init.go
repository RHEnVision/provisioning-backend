//go:build !test

package rbac

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetRbacClient = newClient
}
