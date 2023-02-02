package stubs

import (
	"context"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/models"
)

type AzureClientStub struct {
	createdVms []*armcompute.VirtualMachine
	createdRgs []*armresources.ResourceGroup
}

func DidCreateAzureResourceGroup(ctx context.Context, name string) bool {
	client, err := getAzureClientStub(ctx)
	if err != nil {
		return false
	}
	for _, rg := range client.createdRgs {
		if *rg.Name == name {
			return true
		}
	}
	return false
}

func DidCreateAzureVM(ctx context.Context, name string) bool {
	client, err := getAzureClientStub(ctx)
	if err != nil {
		return false
	}
	for _, vm := range client.createdVms {
		if *vm.Name == name {
			return true
		}
	}
	return false
}

func (stub *AzureClientStub) Status(ctx context.Context) error {
	return nil
}

func (stub *AzureClientStub) CreateVM(ctx context.Context, location string, resourceGroupName string, imageID string, pubkey *models.Pubkey, instanceType clients.InstanceTypeName, vmName string) (*string, error) {
	id := strconv.Itoa(len(stub.createdVms) + 1)

	vm := armcompute.VirtualMachine{
		ID:       &id,
		Name:     &vmName,
		Location: &location,
	}
	stub.createdVms = append(stub.createdVms, &vm)
	return &id, nil
}

func (stub *AzureClientStub) EnsureResourceGroup(ctx context.Context, name string, location string) (*string, error) {
	id := strconv.Itoa(len(stub.createdRgs) + 1)

	rg := armresources.ResourceGroup{
		ID:       &id,
		Name:     &name,
		Location: &location,
	}
	stub.createdRgs = append(stub.createdRgs, &rg)
	return &id, nil
}
