package ec2

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/rs/zerolog"
)

type ec2Client struct {
	ec2     *ec2.Client
	assumed bool
}

func init() {
	clients.GetCustomerEC2Client = newAssumedEC2ClientWithRegion
	clients.GetServiceEC2Client = newEC2ClientWithRegion
}

func logger(ctx context.Context) zerolog.Logger {
	return ctxval.Logger(ctx).With().Str("client", "ec2").Logger()
}

func endpointResolver() aws.EndpointResolverWithOptionsFunc {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           fmt.Sprintf("https://%s.%s.amazonaws.com", service, config.AWS.SigningRegion),
			SigningRegion: config.AWS.SigningRegion,
		}, nil
	})
}

func newEC2ClientWithRegion(ctx context.Context, region string) (clients.EC2, error) {
	if region == "" {
		region = config.AWS.DefaultRegion
	}

	newCfg, err := awsCfg.LoadDefaultConfig(ctx,
		awsCfg.WithRegion(region),
		awsCfg.WithEndpointResolverWithOptions(endpointResolver()),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AWS.Key, config.AWS.Secret, config.AWS.Session)))
	if err != nil {
		return nil, fmt.Errorf("cannot create a new ec2 config: %w", err)
	}

	return &ec2Client{
		ec2:     ec2.NewFromConfig(newCfg),
		assumed: false,
	}, nil
}

func newAssumedEC2ClientWithRegion(ctx context.Context, auth *clients.Authentication, region string) (clients.EC2, error) {
	if typeErr := auth.MustBe(models.ProviderTypeAWS); typeErr != nil {
		return nil, fmt.Errorf("unexpected authentication: %w", typeErr)
	}

	if region == "" {
		region = config.AWS.DefaultRegion
	}

	creds, err := getStsAssumedCredentials(ctx, auth.Payload, region)
	if err != nil {
		return nil, err
	}

	newCfg, err := awsCfg.LoadDefaultConfig(ctx,
		awsCfg.WithRegion(region),
		awsCfg.WithEndpointResolverWithOptions(endpointResolver()),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(*creds.AccessKeyId, *creds.SecretAccessKey, *creds.SessionToken)))
	if err != nil {
		return nil, fmt.Errorf("cannot create a new ec2 config: %w", err)
	}

	return &ec2Client{
		ec2:     ec2.NewFromConfig(newCfg),
		assumed: true,
	}, nil
}

func (c *ec2Client) Status(ctx context.Context) error {
	_, err := c.ListAllRegions(ctx)
	return err
}

func getStsAssumedCredentials(ctx context.Context, arn string, region string) (*stsTypes.Credentials, error) {
	logger := logger(ctx)

	cfg, err := awsCfg.LoadDefaultConfig(ctx,
		awsCfg.WithRegion(region),
		awsCfg.WithEndpointResolverWithOptions(endpointResolver()),
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AWS.Key, config.AWS.Secret, config.AWS.Session)))
	if err != nil {
		logger.Error().Err(err).Msgf("Cannot create an sts client %s", err)
		return nil, fmt.Errorf("cannot create an sts client %w", err)
	}
	stsClient := sts.NewFromConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create an sts client %w", err)
	}

	output, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         ptr.To(arn),
		RoleSessionName: ptr.To("name"),
	})
	if err != nil {
		logger.Error().Err(err).Msgf("cannot assume role %s", err)
		return nil, fmt.Errorf("cannot assume role %w", err)
	}

	return output.Credentials, nil
}

// ImportPubkey imports a key and returns AWS ID
func (c *ec2Client) ImportPubkey(ctx context.Context, key *models.Pubkey, tag string) (string, error) {
	if !c.assumed {
		return "", http.ServiceAccountUnsupportedOperationErr
	}
	logger := logger(ctx)
	logger.Trace().Msgf("Importing AWS key-pair named '%s' with tag '%s'", key.Name, tag)

	input := &ec2.ImportKeyPairInput{}
	input.KeyName = ptr.To(key.Name)
	input.PublicKeyMaterial = []byte(key.Body)
	input.TagSpecifications = []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeKeyPair,
			Tags: []types.Tag{
				{
					Key:   ptr.To("rhhc:id"),
					Value: ptr.To(tag),
				},
			},
		},
	}
	output, err := c.ec2.ImportKeyPair(ctx, input)
	if err != nil {
		if isAWSUnauthorizedError(err) {
			err = clients.UnauthorizedErr
		} else if isAWSOperationError(err, "InvalidKeyPair.Duplicate") {
			err = http.DuplicatePubkeyErr
		}
		return "", fmt.Errorf("cannot import SSH key %s: %w", key.Name, err)
	}

	return ptr.From(output.KeyPairId), nil
}

func (c *ec2Client) DeleteSSHKey(ctx context.Context, handle string) error {
	if !c.assumed {
		return http.ServiceAccountUnsupportedOperationErr
	}
	logger := logger(ctx)
	logger.Trace().Msgf("Deleting AWS key-pair with handle %s", handle)

	input := &ec2.DeleteKeyPairInput{}
	input.KeyPairId = ptr.To(handle)
	_, err := c.ec2.DeleteKeyPair(ctx, input)
	if err != nil {
		if isAWSUnauthorizedError(err) {
			err = clients.UnauthorizedErr
		}
		return fmt.Errorf("cannot delete SSH key %v: %w", input.KeyPairId, err)
	}

	return nil
}

func (c *ec2Client) ListAllRegions(ctx context.Context) ([]clients.Region, error) {
	input := &ec2.DescribeRegionsInput{
		AllRegions: ptr.To(true),
	}

	output, err := c.ec2.DescribeRegions(ctx, input)
	if err != nil {
		if isAWSUnauthorizedError(err) {
			err = clients.UnauthorizedErr
		}
		return nil, fmt.Errorf("cannot list regions: %w", err)
	}

	result := make([]clients.Region, 0, len(output.Regions))
	for _, region := range output.Regions {
		result = append(result, clients.Region(*region.RegionName))
	}

	return result, nil
}

func (c *ec2Client) ListAllZones(ctx context.Context, region clients.Region) ([]clients.Zone, error) {
	input := &ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: ptr.To(true),
		Filters: []types.Filter{
			{
				Name:   ptr.To("region-name"),
				Values: []string{region.String()},
			},
		},
	}

	output, err := c.ec2.DescribeAvailabilityZones(ctx, input)
	if err != nil {
		if isAWSUnauthorizedError(err) {
			err = clients.UnauthorizedErr
		}
		return nil, fmt.Errorf("cannot list zones: %w", err)
	}

	result := make([]clients.Zone, 0, len(output.AvailabilityZones))
	for _, zone := range output.AvailabilityZones {
		result = append(result, clients.Zone(*zone.ZoneName))
	}

	return result, nil
}

func (c *ec2Client) ListInstanceTypesWithPaginator(ctx context.Context) ([]*clients.InstanceType, error) {
	input := &ec2.DescribeInstanceTypesInput{MaxResults: ptr.ToInt32(100)}
	pag := ec2.NewDescribeInstanceTypesPaginator(c.ec2, input)

	res := make([]types.InstanceTypeInfo, 0, 128)
	for pag.HasMorePages() {
		resp, err := pag.NextPage(ctx)
		if err != nil {
			if isAWSUnauthorizedError(err) {
				err = clients.UnauthorizedErr
			}
			return nil, fmt.Errorf("cannot list instance types: %w", err)
		}
		res = append(res, resp.InstanceTypes...)
	}

	// convert to the client type
	instances, err := NewInstanceTypes(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("cannot convert instance types: %w", err)
	}

	return instances, nil
}

func (c *ec2Client) RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error) {
	if !c.assumed {
		return nil, nil, http.ServiceAccountUnsupportedOperationErr
	}
	logger := logger(ctx)
	logger.Trace().Msg("Run AWS EC2 instance")

	encodedUserData := base64.StdEncoding.EncodeToString(userData)
	input := &ec2.RunInstancesInput{
		MaxCount:     ptr.To(amount),
		MinCount:     ptr.To(amount),
		InstanceType: instanceType,
		ImageId:      ptr.To(AMI),
		KeyName:      &keyName,
		UserData:     &encodedUserData,
	}
	if name != nil {
		input.TagSpecifications = []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   ptr.To("Name"),
						Value: name,
					},
				},
			},
		}
	}
	resp, err := c.ec2.RunInstances(ctx, input)
	if err != nil {
		if isAWSUnauthorizedError(err) {
			err = clients.UnauthorizedErr
		}
		return nil, nil, fmt.Errorf("cannot run instances: %w", err)
	}
	instances := c.parseRunInstancesResponse(resp)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ParseRunInstancesResponse: %w", err)
	}
	return instances, resp.ReservationId, nil
}

func (c *ec2Client) parseRunInstancesResponse(respAWS *ec2.RunInstancesOutput) []*string {
	instances := respAWS.Instances
	list := make([]*string, len(instances))
	for i, instance := range instances {
		list[i] = instance.InstanceId
	}
	return list
}
