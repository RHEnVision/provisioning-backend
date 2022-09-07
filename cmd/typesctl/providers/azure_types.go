package providers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/azure/types"
	"github.com/RHEnVision/provisioning-backend/internal/config"
)

func init() {
	provider := TypeProvider{
		PrintRegisteredTypes:      printRegisteredTypesAzure,
		PrintRegionalAvailability: printRegionalAvailabilityAzure,
		GenerateTypes:             generateTypesAzure,
	}
	TypeProviders["azure"] = provider
}

func printRegisteredTypesAzure(name string) {
	types.PrintRegisteredTypes(name)
}

func printRegionalAvailabilityAzure(region, zone string) {
	types.PrintRegionalAvailability(region, zone)
}

func generateTypesAzure() error {
	azureInstanceTypes := clients.NewRegisteredInstanceTypes()
	azureRegionalTypes := clients.NewRegionalInstanceTypes()
	restricted := make(map[armcompute.ResourceSKURestrictionsReasonCode]int, 0)

	opts := azidentity.ClientSecretCredentialOptions{}
	identityClient, err := azidentity.NewClientSecretCredential(config.Azure.TenantID, config.Azure.ClientID, config.Azure.ClientSecret, &opts)
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}
	ctx := context.Background()

	skuClient, err := armcompute.NewResourceSKUsClient(config.Azure.SubscriptionID, identityClient, nil)
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}
	skuPager := skuClient.NewListPager(nil)
	for skuPager.More() {
		nextResult, pageErr := skuPager.NextPage(ctx)
		if pageErr != nil {
			return fmt.Errorf("unable to generate types: %w", pageErr)
		}
		for _, v := range nextResult.Value {
			if *v.ResourceType != "virtualMachines" {
				continue
			}
			instanceType := clients.InstanceType{
				Name:        clients.InstanceTypeName(*v.Name),
				AzureDetail: &clients.InstanceTypeDetailAzure{},
			}
			vcpusPerCore := int32(0)
			for _, c := range v.Capabilities {
				if *c.Name == "CpuArchitectureType" {
					instanceType.Architecture, pageErr = clients.MapArchitectures(ctx, *c.Value)
					if pageErr != nil {
						return fmt.Errorf("unable to generate types: %w", pageErr)
					}
				}
				if *c.Name == "vCPUs" {
					vcpus, vcpuErr := strconv.ParseInt(*c.Value, 10, 32)
					if vcpuErr != nil {
						return fmt.Errorf("unable to generate types: %w", vcpuErr)
					}
					instanceType.VCPUs = int32(vcpus)
				}
				if *c.Name == "vCPUsPerCore" {
					value, vcpupcErr := strconv.ParseInt(*c.Value, 10, 32)
					if vcpupcErr != nil {
						return fmt.Errorf("unable to generate types: %w", vcpupcErr)
					}
					vcpusPerCore = int32(value)
				}
				if *c.Name == "MaxResourceVolumeMB" {
					mbs, volErr := strconv.Atoi(*c.Value)
					if volErr != nil {
						return fmt.Errorf("unable to generate types: %w", volErr)
					}
					instanceType.SetEphemeralStorageFromMB(int64(mbs))
				}
				// Appears to be GiB not GB.
				if *c.Name == "MemoryGB" {
					memoryGB, memErr := strconv.ParseFloat(*c.Value, 64)
					if memErr != nil {
						return fmt.Errorf("unable to generate types: %w", memErr)
					}
					instanceType.MemoryMiB = int64(memoryGB * 1000)
				}
				if *c.Name == "HyperVGenerations" {
					if strings.Contains(*c.Value, "V1") {
						instanceType.AzureDetail.GenV1 = true
					}
					if strings.Contains(*c.Value, "V2") {
						instanceType.AzureDetail.GenV2 = true
					}
				}
			}
			if vcpusPerCore != 0 {
				instanceType.Cores = instanceType.VCPUs / vcpusPerCore
			} else {
				// Some types have no HT (HPC AMD CPUs) thus vcpus per core is zero:
				// Standard_B16ms, Standard_B20ms, Standard_HB60rs, Standard_HC44rs,
				// Standard_M208ms_v2, Standard_M208s_v2, Standard_M416ms_v2,
				// Standard_M416s_v2 and Standard_PB6s
				instanceType.Cores = instanceType.VCPUs
			}

			// Register instance type
			azureInstanceTypes.Register(instanceType)

			if v.Restrictions == nil || len(v.Restrictions) == 0 {
				// Unrestricted type
				for _, location := range v.LocationInfo {
					for _, zone := range location.Zones {
						azureRegionalTypes.Add(strings.ToLower(*location.Location), strings.ToLower(*zone), instanceType)
					}
				}
			} else {
				// Restrictions as documented on Azure: Quota Id is set when the SKU has requiredQuotas parameter
				// as the subscription does not belong to that quota. The "NotAvailableForSubscription" is related
				// to capacity at DC. Possible values include: 'QuotaId', 'NotAvailableForSubscription'
				for _, r := range v.Restrictions {
					if _, ok := restricted[*r.ReasonCode]; !ok {
						restricted[*r.ReasonCode] = 1
					} else {
						restricted[*r.ReasonCode] += 1
					}
				}
			}
		}
	}

	err = azureInstanceTypes.Save("internal/clients/http/azure/types/types.yaml")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = azureRegionalTypes.Save("internal/clients/http/azure/types/availability")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	for key, value := range restricted {
		println("number of", key, "restrictions:", value)
		if key == armcompute.ResourceSKURestrictionsReasonCodeQuotaID {
			println("Increase account quota in Subscription - Documentation - Usage and quotas")
			println("to avoid instance types being restricted from the SKU list.")
		}
	}

	return nil
}
