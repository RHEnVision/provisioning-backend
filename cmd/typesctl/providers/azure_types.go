package providers

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/preload"
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
	preload.AzureInstanceType.PrintRegisteredTypes(name)
}

func printRegionalAvailabilityAzure(region, zone string) {
	preload.AzureInstanceType.PrintRegionalAvailability(region, zone)
}

func generateTypesAzure() error {
	instanceTypes := clients.NewRegisteredInstanceTypes()
	regionalTypes := clients.NewRegionalInstanceTypes()

	ctx := context.Background()
	sc, err := clients.GetServiceAzureClient(ctx)
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = sc.RegisterInstanceTypes(ctx, instanceTypes, regionalTypes)
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = instanceTypes.Save("internal/preload/azure_types.yaml")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = regionalTypes.Save("internal/preload/azure_availability")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	return nil
}
