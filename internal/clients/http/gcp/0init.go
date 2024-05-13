//go:build !test

package gcp

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetGCPClient = newGCPClient
	clients.GetServiceGCPClient = newServiceGCPClient
}
