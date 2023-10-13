package gcp

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/RHEnVision/provisioning-backend/internal/identity"
	"github.com/RHEnVision/provisioning-backend/internal/telemetry"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/page"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type gcpClient struct {
	auth    *clients.Authentication
	options []option.ClientOption
}

func init() {
	clients.GetGCPClient = newGCPClient
}

// GCP SDK does not provide a single client, so only configuration can be shared and
// clients need to be created and closed in each function.
// The difference between the customer and service authentication is which Project ID was given: the service or the customer
func newGCPClient(ctx context.Context, auth *clients.Authentication) (clients.GCP, error) {
	options := []option.ClientOption{
		option.WithCredentialsJSON([]byte(config.GCP.JSON)),
		option.WithQuotaProject(auth.Payload),
		option.WithRequestReason(logging.TraceId(ctx)),
	}
	return &gcpClient{
		auth:    auth,
		options: options,
	}, nil
}

func (c *gcpClient) Status(ctx context.Context) error {
	_, err := c.ListAllRegions(ctx)
	return err
}

func (c *gcpClient) ListAllRegions(ctx context.Context) ([]clients.Region, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListAllRegions")
	defer span.End()

	client, err := compute.NewRegionsRESTClient(ctx, c.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCP regions client: %w", err)
	}
	defer client.Close()

	// This request returns approx. 6kB response (gzipped). Although REST API allows allow-list
	// of fields via the 'fields' URL param ("items.name,items.zones"), gRPC does not allow that.
	// Therefore, we must download all the information only to extract region and zone names.
	req := &computepb.ListRegionsRequest{
		Project: c.auth.Payload,
	}
	iter := client.List(ctx, req)
	regions := make([]clients.Region, 0, 32)
	for {
		region, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			return nil, fmt.Errorf("iterator error: %w", err)
		}
		regions = append(regions, clients.Region(*region.Name))
	}
	return regions, nil
}

func (c *gcpClient) newInstancesClient(ctx context.Context) (*compute.InstancesClient, error) {
	client, err := compute.NewInstancesRESTClient(ctx, c.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCP regions client: %w", err)
	}
	return client, nil
}

func (c *gcpClient) NewInstanceTemplatesClient(ctx context.Context) (*compute.InstanceTemplatesClient, error) {
	client, err := compute.NewInstanceTemplatesRESTClient(ctx, c.options...)
	if err != nil {
		return nil, fmt.Errorf("unable to create GCP templates client: %w", err)
	}
	return client, nil
}

func (c *gcpClient) ListLaunchTemplates(ctx context.Context) ([]*clients.LaunchTemplate, string, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListLaunchTemplates")
	defer span.End()
	var token string
	logger := logger(ctx)
	logger.Trace().Msgf("Listing launch templates")

	templatesClient, err := c.NewInstanceTemplatesClient(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get instances client")
		return nil, "", fmt.Errorf("unable to get instances client: %w", err)
	}
	defer templatesClient.Close()

	limit := page.Limit(ctx).Int()
	req := &computepb.ListInstanceTemplatesRequest{
		Project:    c.auth.Payload,
		MaxResults: ptr.To(uint32(limit)),
		PageToken:  ptr.To(page.Token(ctx)),
	}

	var lst []*computepb.InstanceTemplate
	it := templatesClient.List(ctx, req)
	pager := iterator.NewPager(it, it.PageInfo().MaxSize, token)
	nextToken, err := pager.NextPage(&lst)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get next page: %w", err)
	}

	templatesList := make([]*clients.LaunchTemplate, 0, len(lst))
	for _, template := range lst {
		id := strconv.FormatUint(*template.Id, 10)
		templatesList = append(templatesList, &clients.LaunchTemplate{ID: id, Name: template.GetName()})
	}

	return templatesList, nextToken, nil
}

func (c *gcpClient) InsertInstances(ctx context.Context, params *clients.GCPInstanceParams, amount int64) ([]*string, *string, error) {
	ctx, span := telemetry.StartSpan(ctx, "InsertInstances")
	defer span.End()

	logger := logger(ctx)
	logger.Trace().Msgf("Executing bulk insert for name: %s", *params.NamePattern)

	client, err := c.newInstancesClient(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get instances client")
		return nil, nil, fmt.Errorf("unable to get instances client: %w", err)
	}
	defer client.Close()

	if params.Zone == "" {
		params.Zone = config.GCP.DefaultZone
	}

	pk := models.Pubkey{Body: params.KeyBody}
	pkBody, err := pk.BodyWithUsername(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to get pubkey body with username: %w", err)
	}

	metadata := []*computepb.Items{
		{
			Key:   ptr.To("ssh-keys"),
			Value: ptr.To(pkBody),
		},
	}
	if params.StartupScript != "" {
		metadata = append(metadata, &computepb.Items{
			Key:   ptr.To("startup-script"),
			Value: ptr.To(params.StartupScript),
		})
	}

	req := &computepb.BulkInsertInstanceRequest{
		Project: c.auth.Payload,
		Zone:    params.Zone,
		BulkInsertInstanceResourceResource: &computepb.BulkInsertInstanceResource{
			NamePattern: params.NamePattern,
			Count:       &amount,
			MinCount:    &amount,
			InstanceProperties: &computepb.InstanceProperties{
				Labels: map[string]string{
					"rh-rid":  config.EnvironmentPrefix("r", strconv.FormatInt(params.ReservationID, 10)),
					"rh-uuid": params.UUID,
					"rh-org":  identity.Identity(ctx).Identity.OrgID,
				},
				NetworkInterfaces: []*computepb.NetworkInterface{
					{
						AccessConfigs: []*computepb.AccessConfig{
							{
								Name: ptr.To("External NAT"),
								Type: ptr.To("ONE_TO_ONE_NAT"),
							},
						},
						Name: ptr.To("global/networks/default"),
					},
				},
				Metadata: &computepb.Metadata{
					Items: metadata,
				},
			},
		},
	}

	if params.LaunchTemplateID != "" {
		template := fmt.Sprintf("global/instanceTemplates/%s", params.LaunchTemplateID)
		req.BulkInsertInstanceResourceResource.SourceInstanceTemplate = &template
	}

	if params.MachineType != "" {
		req.BulkInsertInstanceResourceResource.InstanceProperties.MachineType = ptr.To(params.MachineType)
	}

	if params.ImageName != "" {
		req.BulkInsertInstanceResourceResource.InstanceProperties.Disks = []*computepb.AttachedDisk{
			{
				InitializeParams: &computepb.AttachedDiskInitializeParams{
					SourceImage: &params.ImageName,
				},
				AutoDelete: ptr.To(true),
				Boot:       ptr.To(true),
				Type:       ptr.To(computepb.AttachedDisk_PERSISTENT.String()),
			},
		}
	}

	op, err := client.BulkInsert(ctx, req)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		logger.Error().Err(err).Msg("Bulk insert operation failed")
		return nil, nil, fmt.Errorf("cannot bulk insert instances: %w", err)
	}
	if err = op.Wait(ctx); err != nil {
		logger.Error().Err(err).Msg("Bulk wait operation failed")
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, fmt.Errorf("cannot bulk insert instances: %w", err)
	}

	if !op.Done() {
		return nil, ptr.To(op.Name()), fmt.Errorf("an error occured on operation %s: %w", op.Name(), ErrOperationFailed)
	}

	ids, err := c.ListInstancesIDsByLabel(ctx, params.UUID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot list instances ids: %w", err)
	}
	return ids, ptr.To(op.Name()), nil
}

func (c *gcpClient) ListInstancesIDsByLabel(ctx context.Context, uuid string) ([]*string, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListInstancesIDsByLabel")
	defer span.End()

	logger := logger(ctx)
	ids := make([]*string, 0)
	filter := fmt.Sprintf("labels.rh-uuid=%v", uuid)
	lstReq := &computepb.AggregatedListInstancesRequest{
		Project: c.auth.Payload,
		Filter:  &filter,
	}

	client, err := c.newInstancesClient(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get instances client")
		return nil, fmt.Errorf("unable to get instances client: %w", err)
	}
	defer client.Close()

	instances := client.AggregatedList(ctx, lstReq)
	logger.Trace().Msg("Fetching instance ids")
	for {
		pair, err := instances.Next()
		if errors.Is(err, iterator.Done) {
			break
		} else if err != nil {
			logger.Error().Err(err).Msg("An error occurred during fetching instance ids")
			span.SetStatus(codes.Error, err.Error())
			return nil, fmt.Errorf("cannot fetch instance ids: %w", err)
		} else {
			instances := pair.Value.Instances
			for _, insta := range instances {
				idAsString := strconv.FormatUint(insta.GetId(), 10)
				ids = append(ids, &idAsString)
			}
		}
	}
	return ids, nil
}

func (c *gcpClient) GetInstanceDescriptionByID(ctx context.Context, id, zone string) (*clients.InstanceDescription, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetInstanceDescriptionByID")
	defer span.End()

	logger := logger(ctx)

	client, err := c.newInstancesClient(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Could not get instances client")
		return nil, fmt.Errorf("unable to get instances client: %w", err)
	}
	defer client.Close()

	projectId := c.auth.String()

	instance, err := client.Get(ctx, &computepb.GetInstanceRequest{Instance: id, Project: projectId, Zone: zone})
	if err != nil {
		return nil, fmt.Errorf("unable to get instance: %w", err)
	}
	instanceId := strconv.FormatUint(instance.GetId(), 10)
	instanceDesc := clients.InstanceDescription{ID: instanceId}
	for _, n := range instance.NetworkInterfaces {
		instanceDesc.PrivateIPv4 = ptr.FromOrEmpty(n.NetworkIP)
		if len(n.AccessConfigs) > 0 && n.AccessConfigs[0] != nil {
			instanceDesc.IPv4 = *n.AccessConfigs[0].NatIP
			break
		}
	}
	return &instanceDesc, nil
}
