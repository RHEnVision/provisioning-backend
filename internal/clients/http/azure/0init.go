//go:build !test

package azure

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetAzureClient = newAzureClient
	clients.GetServiceAzureClient = newServiceClient
}
