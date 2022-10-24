package providers

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/azure/types"
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

	err = instanceTypes.Save("internal/clients/http/azure/types/types.yaml")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = regionalTypes.Save("internal/clients/http/azure/types/availability")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	return nil
}
