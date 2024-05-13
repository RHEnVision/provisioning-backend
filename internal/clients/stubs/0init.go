//go:build test

package stubs

import "github.com/RHEnVision/provisioning-backend/internal/clients"

func init() {
	clients.GetAzureClient = getAzureClient
	clients.GetEC2Client = newEC2CustomerClientStubWithRegion
	clients.GetServiceEC2Client = newEC2ServiceClientStubWithRegion
	clients.GetGCPClient = newGCPCustomerClientStub
	clients.GetServiceGCPClient = getServiceGCPClientStub
	clients.GetImageBuilderClient = getImageBuilderClientStub
	clients.GetRbacClient = getRbacClient
	clients.GetSourcesClient = getSourcesClient
}
