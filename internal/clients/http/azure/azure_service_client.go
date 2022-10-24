package azure

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
)

type serviceClient struct {
	credential *azidentity.ClientSecretCredential
}

func init() {
	clients.GetServiceAzureClient = newServiceClient
}

func newServiceClient(ctx context.Context) (clients.ServiceAzure, error) {
	opts := azidentity.ClientSecretCredentialOptions{}
	identityClient, err := azidentity.NewClientSecretCredential(config.Azure.TenantID, config.Azure.ClientID, config.Azure.ClientSecret, &opts)
	if err != nil {
		return nil, fmt.Errorf("unable to init Azure credentials: %w", err)
	}

	return &serviceClient{
		credential: identityClient,
	}, nil
}

func (c *serviceClient) RegisterInstanceTypes(ctx context.Context, instanceTypes *clients.RegisteredInstanceTypes, regionalTypes *clients.RegionalTypeAvailability) error {
	restricted := make(map[armcompute.ResourceSKURestrictionsReasonCode]int, 0)

	skuClient, err := armcompute.NewResourceSKUsClient(config.Azure.SubscriptionID, c.credential, nil)
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}
	skuPager := skuClient.NewListPager(nil)
	for skuPager.More() {
		nextResult, pageErr := skuPager.NextPage(ctx)
		if pageErr != nil {
			return fmt.Errorf("unable to generate types: %w", pageErr)
		}
		for _, resourceSKU := range nextResult.Value {
			if *resourceSKU.ResourceType != "virtualMachines" {
				continue
			}
			instanceType, err2 := c.typeFromSKU(ctx, resourceSKU)
			if err2 != nil {
				return err2
			}

			// Register instance type
			instanceTypes.Register(instanceType)

			if resourceSKU.Restrictions == nil || len(resourceSKU.Restrictions) == 0 {
				// Unrestricted type
				for _, location := range resourceSKU.LocationInfo {
					for _, zone := range location.Zones {
						regionalTypes.Add(strings.ToLower(*location.Location), strings.ToLower(*zone), instanceType)
					}
				}
			} else {
				// Restrictions as documented on Azure: Quota Id is set when the SKU has requiredQuotas parameter
				// as the subscription does not belong to that quota. The "NotAvailableForSubscription" is related
				// to capacity at DC. Possible values include: 'QuotaId', 'NotAvailableForSubscription'
				for _, r := range resourceSKU.Restrictions {
					if _, ok := restricted[*r.ReasonCode]; !ok {
						restricted[*r.ReasonCode] = 1
					} else {
						restricted[*r.ReasonCode] += 1
					}
				}
			}
		}
	}

	for key, value := range restricted {
		logger := logger(ctx)
		logger.Trace().Msgf("Number of %s restrictions: %d", key, value)
		if key == armcompute.ResourceSKURestrictionsReasonCodeQuotaID {
			logger.Trace().Msg("Increase account quota in Subscription - Documentation - Usage and quotas")
			logger.Trace().Msg("to avoid instance types being restricted from the SKU list.")
		}
	}

	return nil
}

func (c *serviceClient) typeFromSKU(ctx context.Context, v *armcompute.ResourceSKU) (clients.InstanceType, error) {
	var err error
	instanceType := clients.InstanceType{
		Name:        clients.InstanceTypeName(*v.Name),
		AzureDetail: &clients.InstanceTypeDetailAzure{},
	}
	vcpusPerCore := int32(0)

	for _, c := range v.Capabilities {
		switch *c.Name {
		case "CpuArchitectureType":
			instanceType.Architecture, err = clients.MapArchitectures(ctx, *c.Value)
			if err != nil {
				return clients.InstanceType{}, fmt.Errorf("unable to generate types: %w", err)
			}
		case "vCPUs":
			vcpus, vcpuErr := strconv.ParseInt(*c.Value, 10, 32)
			if vcpuErr != nil {
				return clients.InstanceType{}, fmt.Errorf("unable to generate types: %w", vcpuErr)
			}
			instanceType.VCPUs = int32(vcpus)
		case "vCPUsPerCore":
			value, vcpupcErr := strconv.ParseInt(*c.Value, 10, 32)
			if vcpupcErr != nil {
				return clients.InstanceType{}, fmt.Errorf("unable to generate types: %w", vcpupcErr)
			}
			vcpusPerCore = int32(value)
		case "MaxResourceVolumeMB":
			mbs, volErr := strconv.Atoi(*c.Value)
			if volErr != nil {
				return clients.InstanceType{}, fmt.Errorf("unable to generate types: %w", volErr)
			}
			instanceType.SetEphemeralStorageFromMB(int64(mbs))
		case "MemoryGB":
			// Appears to be GiB not GB.
			memoryGB, memErr := strconv.ParseFloat(*c.Value, 64)
			if memErr != nil {
				return clients.InstanceType{}, fmt.Errorf("unable to generate types: %w", memErr)
			}
			instanceType.MemoryMiB = int64(memoryGB * 1000)
		case "HyperVGenerations":
			instanceType.AzureDetail.GenV1 = strings.Contains(*c.Value, "V1")
			instanceType.AzureDetail.GenV2 = strings.Contains(*c.Value, "V2")
		}
	}

	instanceType.Cores = instanceType.VCPUs
	// Some types have no HT (HPC AMD CPUs) thus vcpus per core is zero (Standard_PB6s...)
	if vcpusPerCore != 0 {
		instanceType.Cores /= vcpusPerCore
	}

	return instanceType, nil
}
