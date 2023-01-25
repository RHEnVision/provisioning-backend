package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
)

type client struct {
	subscriptionID string
	credential     *azidentity.ClientSecretCredential
}

func init() {
	clients.GetAzureClient = newAzureClient
}

func newAzureClient(ctx context.Context, auth *clients.Authentication) (clients.Azure, error) {
	opts := azidentity.ClientSecretCredentialOptions{}
	identityClient, err := azidentity.NewClientSecretCredential(config.Azure.TenantID, config.Azure.ClientID, config.Azure.ClientSecret, &opts)
	if err != nil {
		return nil, fmt.Errorf("unable to init Azure credentials: %w", err)
	}

	return &client{
		subscriptionID: auth.Payload,
		credential:     identityClient,
	}, nil
}

func (c *client) newResourceGroupsClient(ctx context.Context) (*armresources.ResourceGroupsClient, error) {
	client, err := armresources.NewResourceGroupsClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create resources Azure client: %w", err)
	}
	return client, nil
}

func (c *client) newVirtualMachinesClient(ctx context.Context) (*armcompute.VirtualMachinesClient, error) {
	vmClient, err := armcompute.NewVirtualMachinesClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create VM Azure client: %w", err)
	}
	return vmClient, nil
}

func (c *client) newSubscriptionsClient(ctx context.Context) (*armsubscriptions.Client, error) {
	client, err := armsubscriptions.NewClient(c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create subscriptioons Azure client: %w", err)
	}
	return client, nil
}

func (c *client) newSshKeysClient(ctx context.Context) (*armcompute.SSHPublicKeysClient, error) {
	client, err := armcompute.NewSSHPublicKeysClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create SSH keys Azure client: %w", err)
	}
	return client, nil
}

func (c *client) newVirtualNetworksClient(ctx context.Context) (*armnetwork.VirtualNetworksClient, error) {
	vnetClient, err := armnetwork.NewVirtualNetworksClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create Virtual networks Azure client: %w", err)
	}
	return vnetClient, nil
}

func (c *client) newSubnetsClient(ctx context.Context) (*armnetwork.SubnetsClient, error) {
	subnetClient, err := armnetwork.NewSubnetsClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create SSH keys Azure client: %w", err)
	}
	return subnetClient, nil
}

func (c *client) newPublicIPAddressesClient(ctx context.Context) (*armnetwork.PublicIPAddressesClient, error) {
	publicIPAddressClient, err := armnetwork.NewPublicIPAddressesClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create public IP addresses Azure client: %w", err)
	}
	return publicIPAddressClient, nil
}

func (c *client) newSecurityGroupsClient(ctx context.Context) (*armnetwork.SecurityGroupsClient, error) {
	nsgClient, err := armnetwork.NewSecurityGroupsClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create security groups Azure client: %w", err)
	}
	return nsgClient, nil
}

func (c *client) newInterfacesClient(ctx context.Context) (*armnetwork.InterfacesClient, error) {
	nicClient, err := armnetwork.NewInterfacesClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create interfaces Azure client: %w", err)
	}
	return nicClient, nil
}

func (c *client) Status(ctx context.Context) error {
	client, err := c.newSubscriptionsClient(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize status request: %w", err)
	}
	_, err = client.Get(ctx, c.subscriptionID, nil)
	if err != nil {
		return fmt.Errorf("unable to perform status request: %w", err)
	}
	return nil
}
