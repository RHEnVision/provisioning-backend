package gcp

import (
	"context"
	"fmt"

	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/aws/smithy-go/ptr"
	"github.com/rs/zerolog"
)

type gcpClient struct {
	instancesClient *compute.InstancesClient
	context         context.Context
	logger          zerolog.Logger
}

func init() {
	clients.GetGCPClient = newGCPClient
}

func newGCPClient(ctx context.Context) (clients.GCP, error) {
	logger := ctxval.Logger(ctx).With().Str("client", "gcp").Logger()

	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create new instances REST client: %w", err)
	}

	return &gcpClient{
		instancesClient: instancesClient,
		context:         ctx,
		logger:          logger,
	}, nil
}

func (c *gcpClient) Close() {
	c.instancesClient.Close()
}

func (c *gcpClient) RunInstances(ctx context.Context, projectID string, namePattern *string, imageName *string, amount int64, machineType string, zone string, keyBody string) error {
	req := &computepb.BulkInsertInstanceRequest{
		Project: projectID,
		Zone:    zone,
		BulkInsertInstanceResourceResource: &computepb.BulkInsertInstanceResource{
			NamePattern: namePattern,
			Count:       ptr.Int64(amount),
			MinCount:    ptr.Int64(amount),
			InstanceProperties: &computepb.InstanceProperties{
				Disks: []*computepb.AttachedDisk{
					{
						InitializeParams: &computepb.AttachedDiskInitializeParams{
							SourceImage: imageName,
						},
						AutoDelete: ptr.Bool(true),
						Boot:       ptr.Bool(true),
						Type:       ptr.String(computepb.AttachedDisk_PERSISTENT.String()),
					},
				},
				MachineType: ptr.String(machineType),
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						Name: ptr.String("global/networks/default"),
					},
				},
				Metadata: &computepb.Metadata{
					Items: []*computepb.Items{
						{
							Key:   ptr.String("ssh-keys"),
							Value: ptr.String(keyBody),
						},
					},
				},
			},
		},
	}

	c.logger.Trace().Msg("Executing Insert")
	_, err := c.instancesClient.BulkInsert(ctx, req)
	if err != nil {
		return fmt.Errorf("cannot bulk insert instances: %w", err)
	}

	return nil
}
