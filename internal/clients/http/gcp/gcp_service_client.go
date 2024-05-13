package gcp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type gcpServiceClient struct {
	options []option.ClientOption
}

func newServiceGCPClient(ctx context.Context) (clients.ServiceGCP, error) {
	options := []option.ClientOption{
		option.WithCredentialsJSON([]byte(config.GCP.JSON)),
		option.WithRequestReason(logging.TraceId(ctx)),
	}
	return &gcpServiceClient{
		options: options,
	}, nil
}

func (c *gcpServiceClient) RegisterInstanceTypes(ctx context.Context, instanceTypes *clients.RegisteredInstanceTypes, regionalTypes *clients.RegionalTypeAvailability) error {
	_, zones, err := c.ListAllRegionsAndZones(ctx)
	if err != nil {
		return fmt.Errorf("unable to list GCP regions and zones: %w", err)
	}

	for _, zone := range zones {
		instTypes, zoneErr := c.ListMachineTypes(ctx, getZone(zone))
		if zoneErr != nil {
			return fmt.Errorf("unable to list gcp machine types: %w", zoneErr)
		}
		for _, instanceType := range instTypes {
			instanceTypes.Register(*instanceType)
			regionalTypes.Add(zone.String(), "", *instanceType)
		}
	}
	return nil
}

func (c *gcpServiceClient) newRegionsClient(ctx context.Context) (*compute.RegionsClient, error) {
	client, err := compute.NewRegionsRESTClient(ctx, c.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCP regions client: %w", err)
	}
	return client, nil
}

func (c *gcpServiceClient) ListAllRegionsAndZones(ctx context.Context) ([]clients.Region, []clients.Zone, error) {
	client, err := c.newRegionsClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to list regions and zones: %w", err)
	}
	defer client.Close()

	// This request returns approx. 6kB response (gzipped). Although REST API allows allow-list
	// of fields via the 'fields' URL param ("items.name,items.zones"), gRPC does not allow that.
	// Therefore, we must download all the information only to extract region and zone names.
	req := &computepb.ListRegionsRequest{
		Project: config.GCP.ProjectID,
	}
	iter := client.List(ctx, req)
	regions := make([]clients.Region, 0, 32)
	zones := make([]clients.Zone, 0, 64)
	for {
		region, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("iterator error: %w", err)
		}
		regions = append(regions, clients.Region(region.GetName()))
		for _, zone := range region.GetZones() {
			zoneAsArr := strings.Split(zone, "/")
			zones = append(zones, clients.Zone(zoneAsArr[len(zoneAsArr)-1]))
		}
	}
	return regions, zones, nil
}

func (c *gcpServiceClient) newMachineTypeClient(ctx context.Context) (*compute.MachineTypesClient, error) {
	client, err := compute.NewMachineTypesRESTClient(ctx, c.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCP regions client: %w", err)
	}
	return client, nil
}

func (c *gcpServiceClient) ListMachineTypes(ctx context.Context, zone string) ([]*clients.InstanceType, error) {
	client, err := c.newMachineTypeClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list machine types: %w", err)
	}
	defer client.Close()

	if zone == "" {
		zone = config.GCP.DefaultZone
	}

	req := &computepb.ListMachineTypesRequest{
		Project: config.GCP.ProjectID,
		Zone:    zone,
	}
	iter := client.List(ctx, req)
	machineTypes := make([]*clients.InstanceType, 0, 32)
	for {
		machineType, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("iterator error: %w", err)
		}
		arch := clients.ArchitectureTypeX86_64
		if getMachineFamily(machineType.GetName()) == "t2a" {
			arch = clients.ArchitectureTypeArm64
		}
		machineTypes = append(machineTypes, &clients.InstanceType{
			Name:               clients.InstanceTypeName(machineType.GetName()),
			VCPUs:              machineType.GetGuestCpus(),
			Cores:              0,
			Architecture:       arch,
			MemoryMiB:          mbToMib(float32(machineType.GetMemoryMb())),
			EphemeralStorageGB: getTotalStorage(machineType.GetScratchDisks()),
		})
	}
	return machineTypes, nil
}

func getZone(zone clients.Zone) string {
	zoneAsArr := strings.Split(zone.String(), "/")
	return zoneAsArr[len(zoneAsArr)-1]
}

func getMachineFamily(machineType string) string {
	machineFamily := strings.Split(machineType, "-")
	return machineFamily[0]
}

func getTotalStorage(disks []*computepb.ScratchDisks) int64 {
	var sum int64 = 0
	for _, disk := range disks {
		sum += int64(disk.GetDiskGb())
	}
	return sum
}

func mbToMib(mb float32) int64 {
	mib := mb * 0.9536
	return int64(mib)
}
