package providers

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp/types"
)

func init() {
	provider := TypeProvider{
		PrintRegisteredTypes:      printRegisteredTypesGCP,
		PrintRegionalAvailability: printRegionalAvailabilityGCP,
		GenerateTypes:             generateTypesGCP,
	}
	TypeProviders["gcp"] = provider
}

func printRegisteredTypesGCP(name string) {
	types.PrintRegisteredTypes(name)
}

func printRegionalAvailabilityGCP(region, zone string) {
	types.PrintRegionalAvailability(region, zone)
}

func generateTypesGCP() error {
	instanceTypes := clients.NewRegisteredInstanceTypes()
	regionalTypes := clients.NewRegionalInstanceTypes()
	ctx := context.Background()

	gcpClient, err := clients.GetServiceGCPClient(ctx)
	if err != nil {
		return fmt.Errorf("unable to get GCP client: %w", err)
	}

	err = gcpClient.RegisterInstanceTypes(ctx, instanceTypes, regionalTypes)
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = instanceTypes.Save("internal/clients/http/gcp/types/types.yaml")
	if err != nil {
		return fmt.Errorf("unable to save types: %w", err)
	}

	err = regionalTypes.Save("internal/clients/http/gcp/types/availability")
	if err != nil {
		return fmt.Errorf("unable to save regional types: %w", err)
	}

	return nil
}
