package stubs

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources/v3"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
)

var ErrNotStartedVM = errors.New("the VM under given resumeToken not started")

type AzureClientStub struct {
	startedVms []*armcompute.VirtualMachine
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

func CountStubAzureVMs(ctx context.Context) int {
	client, err := getAzureClientStub(ctx)
	if err != nil {
		return 0
	}
	return len(client.createdVms)
}

func (stub *AzureClientStub) Status(ctx context.Context) error {
	return nil
}

func (stub *AzureClientStub) CreateVMs(ctx context.Context, vmParams clients.AzureInstanceParams, amount int64, vmNamePrefix string) ([]clients.InstanceDescription, error) {
	vmIds := make([]clients.InstanceDescription, amount)
	resumeTokens := make([]string, amount)
	var i int64
	var err error
	for i = 0; i < amount; i++ {
		vmName := fmt.Sprintf("%s-%d", vmNamePrefix, int64(len(stub.startedVms))+i)
		resumeTokens[i], err = stub.BeginCreateVM(ctx, vmParams, vmName)
		if err != nil {
			return vmIds, err
		}
	}
	for i = 0; i < amount; i++ {
		instanceID, err := stub.WaitForVM(ctx, resumeTokens[i])
		if err != nil {
			return vmIds, err
		}
		vmIds[i].ID = string(instanceID)
		vmIds[i].IPv4 = fmt.Sprintf("198.51.100.%d", i+1)
	}

	return vmIds, nil
}

func (stub *AzureClientStub) BeginCreateVM(ctx context.Context, vmParams clients.AzureInstanceParams, vmName string) (string, error) {
	id := "with-polling-" + strconv.Itoa(len(stub.startedVms)+1)

	vm := armcompute.VirtualMachine{
		ID:       &id,
		Name:     &vmName,
		Location: &vmParams.Location,
	}
	stub.startedVms = append(stub.startedVms, &vm)
	// we use the id as a resume token
	return id, nil
}

func (stub *AzureClientStub) WaitForVM(ctx context.Context, resumeToken string) (clients.AzureInstanceID, error) {
	for i, vm := range stub.startedVms {
		if *vm.ID == resumeToken {
			stub.createdVms = append(stub.createdVms, vm)
			stub.startedVms = append(stub.startedVms[:i], stub.startedVms[i+1:]...)
			return clients.AzureInstanceID(*vm.ID), nil
		}
	}
	return "", ErrNotStartedVM
}

func (stub *AzureClientStub) EnsureResourceGroup(ctx context.Context, name string, location string) (clients.AzureResourceGroup, error) {
	id := strconv.Itoa(len(stub.createdRgs) + 1)

	if name == "existing" {
		return clients.AzureResourceGroup{ID: "/subscriptions/123/resourceGroups/mockedID", Name: "existing", Location: "westeurope"}, nil
	}

	rg := armresources.ResourceGroup{
		ID:       &id,
		Name:     &name,
		Location: &location,
	}
	stub.createdRgs = append(stub.createdRgs, &rg)
	return clients.AzureResourceGroup{ID: id, Name: name, Location: location}, nil
}

func (stub *AzureClientStub) TenantId(ctx context.Context) (clients.AzureTenantId, error) {
	return "4645f0cb-43f5-4586-b2c9-8d5c58577e3e", nil
}

func (stub *AzureClientStub) ListResourceGroups(ctx context.Context) ([]string, error) {
	return []string{"firstGroup", "secondGroup", "test"}, nil
}
