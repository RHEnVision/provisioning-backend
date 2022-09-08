package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/rs/zerolog"
)

type azureClient struct {
	subscriptionID string
	credential     *azidentity.ClientSecretCredential
}

func init() {
	clients.GetAzureClient = newAzureClient
}

func logger(ctx context.Context) zerolog.Logger {
	return ctxval.Logger(ctx).With().Str("client", "azure").Logger()
}

func newAzureClient(ctx context.Context, auth *clients.Authentication) (clients.Azure, error) {
	opts := azidentity.ClientSecretCredentialOptions{}
	identityClient, err := azidentity.NewClientSecretCredential(config.Azure.TenantID, config.Azure.ClientID, config.Azure.ClientSecret, &opts)
	if err != nil {
		return nil, fmt.Errorf("unable to init Azure credentials: %w", err)
	}

	return &azureClient{
		subscriptionID: auth.Payload,
		credential:     identityClient,
	}, nil
}

func (c *azureClient) newResourcesClient(ctx context.Context) (*armresources.Client, error) {
	client, err := armresources.NewClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create resources Azure client: %w", err)
	}
	return client, nil
}

func (c *azureClient) newSubscriptionsClient(ctx context.Context) (*armsubscriptions.Client, error) {
	client, err := armsubscriptions.NewClient(c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create subscriptioons Azure client: %w", err)
	}
	return client, nil
}

func (c *azureClient) newSshKeysClient(ctx context.Context) (*armcompute.SSHPublicKeysClient, error) {
	client, err := armcompute.NewSSHPublicKeysClient(c.subscriptionID, c.credential, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create SSH keys Azure client: %w", err)
	}
	return client, nil
}

func (c *azureClient) Status(ctx context.Context) error {
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
