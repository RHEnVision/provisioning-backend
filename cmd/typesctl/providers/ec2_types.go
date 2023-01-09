package providers

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2/types"
)

func init() {
	provider := TypeProvider{
		PrintRegisteredTypes:      printRegisteredTypesEC2,
		PrintRegionalAvailability: printRegionalAvailabilityEC2,
		GenerateTypes:             generateTypesEC2,
	}
	TypeProviders["ec2"] = provider
}

func printRegisteredTypesEC2(name string) {
	types.PrintRegisteredTypes(name)
}

func printRegionalAvailabilityEC2(region, zone string) {
	types.PrintRegionalAvailability(region, zone)
}

func generateTypesEC2() error {
	instanceTypes := clients.NewRegisteredInstanceTypes()
	regionalTypes := clients.NewRegionalInstanceTypes()
	ctx := context.Background()

	defaultClient, err := clients.GetServiceEC2Client(ctx, "")
	if err != nil {
		return fmt.Errorf("unable to get default EC2 client: %w", err)
	}

	regions, err := defaultClient.ListAllRegions(ctx)
	if err != nil {
		return fmt.Errorf("unable to list EC2 regions: %w", err)
	}

	// This will throw AuthFailure "AWS was not able to validate the provided access credentials" unless all regions
	// are enabled and "Valid in all AWS Regions" STS endpoint is configured for the account.
	//
	// On some accounts, there is a region that fails to generate with this tool, and it is not visible on the
	// account setting page and cannot be enabled. For this reason, the tool skips such regions.
	//
	// For more info:
	// https://aws.amazon.com/premiumsupport/knowledge-center/iam-validate-access-credentials/
	// https://docs.aws.amazon.com/general/latest/gr/rande-manage.html
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html#sts-regions-manage-tokens
	for _, region := range regions {
		fmt.Printf("Generating for region %s", region)
		client, regionErr := clients.GetServiceEC2Client(ctx, region.String())
		if regionErr != nil {
			return fmt.Errorf("unable to get regional EC2 client: %w", regionErr)
		}
		instTypes, regionErr := client.ListInstanceTypesWithPaginator(ctx)
		if regionErr != nil {
			// skip
			fmt.Printf("unable to list EC2 instance types (region STS not enabled?): %s\n", regionErr.Error())
			continue
		}
		for _, instanceType := range instTypes {
			// filter out i386 architecture as AWS types share the name for 32/64 bit Intel
			if ValidArchitectures.MatchString(instanceType.Architecture.String()) {
				instanceTypes.Register(*instanceType)
				regionalTypes.Add(region.String(), "", *instanceType)
			}
		}
	}

	err = instanceTypes.Save("internal/clients/http/ec2/types/types.yaml")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	err = regionalTypes.Save("internal/clients/http/ec2/types/availability")
	if err != nil {
		return fmt.Errorf("unable to generate types: %w", err)
	}

	return nil
}
