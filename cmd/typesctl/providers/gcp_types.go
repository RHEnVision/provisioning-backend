package providers

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/preload"
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
	preload.GCPInstanceType.PrintRegisteredTypes(name)
}

func printRegionalAvailabilityGCP(region, zone string) {
	preload.GCPInstanceType.PrintRegionalAvailability(region, zone)
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

	err = instanceTypes.Save("internal/preload/gcp_types.yaml")
	if err != nil {
		return fmt.Errorf("unable to save types: %w", err)
	}

	err = regionalTypes.Save("internal/preload/gcp_availability")
	if err != nil {
		return fmt.Errorf("unable to save regional types: %w", err)
	}

	return nil
}
