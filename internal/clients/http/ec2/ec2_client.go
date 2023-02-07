package ec2

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"

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
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const TraceName = "github.com/RHEnVision/provisioning-backend/internal/clients/http/ec2"

type ec2Client struct {
	ec2     *ec2.Client
	sts     *sts.Client
	iam     *iam.Client
	assumed bool
}

func init() {
	clients.GetEC2Client = newAssumedEC2ClientWithRegion
	clients.GetServiceEC2Client = newEC2ClientWithRegion
}

func logger(ctx context.Context) *zerolog.Logger {
	logger := ctxval.Logger(ctx).With().Str("client", "ec2").Logger()
	return &logger
}

func awsConfig(ctx context.Context, region string, optFns ...func(*awsCfg.LoadOptions) error) (*aws.Config, error) {
	if region == "" {
		region = config.AWS.DefaultRegion
	}

	loggingOpt := awsCfg.WithClientLogMode(aws.LogRetries)
	if config.AWS.Logging {
		loggingOpt = awsCfg.WithClientLogMode(aws.LogRequestWithBody | aws.LogRetries | aws.LogResponseWithBody | aws.LogSigning)
	}

	optFns = append(optFns, loggingOpt,
		awsCfg.WithLogger(NewEC2Logger(ctx)),
		awsCfg.WithRegion(region))

	newCfg, err := awsCfg.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &newCfg, nil
}

func newEC2ClientWithRegion(ctx context.Context, region string) (clients.EC2, error) {
	if region == "" {
		region = config.AWS.DefaultRegion
	}

	cfg, err := awsConfig(ctx, region,
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AWS.Key, config.AWS.Secret, config.AWS.Session)))
	if err != nil {
		return nil, fmt.Errorf("aws: %w", err)
	}

	return &ec2Client{
		ec2:     ec2.NewFromConfig(*cfg),
		sts:     sts.NewFromConfig(*cfg),
		iam:     iam.NewFromConfig(*cfg),
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

	assumedCredentials, err := getStsAssumedCredentials(ctx, auth.Payload, region)
	if err != nil {
		return nil, err
	}

	cfg, err := awsConfig(ctx, region,
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			*assumedCredentials.AccessKeyId,
			*assumedCredentials.SecretAccessKey,
			*assumedCredentials.SessionToken)))
	if err != nil {
		return nil, fmt.Errorf("aws: %w", err)
	}

	return &ec2Client{
		ec2:     ec2.NewFromConfig(*cfg),
		sts:     sts.NewFromConfig(*cfg),
		iam:     iam.NewFromConfig(*cfg),
		assumed: true,
	}, nil
}

func (c *ec2Client) Status(ctx context.Context) error {
	_, err := c.ListAllRegions(ctx)
	return err
}

func getStsAssumedCredentials(ctx context.Context, arn string, region string) (*stsTypes.Credentials, error) {
	logger := logger(ctx)

	cfg, err := awsConfig(ctx, region,
		awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.AWS.Key, config.AWS.Secret, config.AWS.Session)))
	if err != nil {
		return nil, fmt.Errorf("aws sts: %w", err)
	}
	stsClient := sts.NewFromConfig(*cfg)
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

// ImportPubkey imports a key and returns AWS KeyPair name.
// The AWS name will be set to value of models.Pubkey Name.
func (c *ec2Client) ImportPubkey(ctx context.Context, key *models.Pubkey, tag string) (string, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "ImportPubkey")
	defer span.End()

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
		span.SetStatus(codes.Error, err.Error())
		return "", fmt.Errorf("cannot import SSH key %s: %w", key.Name, err)
	}

	return *output.KeyPairId, nil
}

func (c *ec2Client) GetPubkeyName(ctx context.Context, fingerprint string) (string, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "fetchPubkeyName")
	defer span.End()

	if !c.assumed {
		return "", http.ServiceAccountUnsupportedOperationErr
	}
	logger := logger(ctx)
	logger.Trace().Msgf("Fetching AWS key with fingerprint '%s' to get its name", fingerprint)
	input := &ec2.DescribeKeyPairsInput{}
	input.Filters = []types.Filter{{Name: ptr.To("fingerprint"), Values: []string{fingerprint}}}
	output, err := c.ec2.DescribeKeyPairs(ctx, input)
	if err != nil {
		if isAWSUnauthorizedError(err) {
			err = clients.UnauthorizedErr
		}
		span.SetStatus(codes.Error, err.Error())
		return "", fmt.Errorf("cannot fetch SSH key to update its tag %s: %w", fingerprint, err)
	}

	if len(output.KeyPairs) == 0 {
		span.SetStatus(codes.Error, fmt.Sprintf("no KeyPair with fingerpring (%s) found", fingerprint))
		return "", fmt.Errorf("SSH key not found by its fingerprint (%s): %w", fingerprint, http.PubkeyNotFoundErr)
	}
	return *output.KeyPairs[0].KeyName, nil
}

func (c *ec2Client) DeleteSSHKey(ctx context.Context, handle string) error {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "DeleteSSHKey")
	defer span.End()

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
		span.SetStatus(codes.Error, err.Error())
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

func (c *ec2Client) ListInstanceTypes(ctx context.Context) ([]*clients.InstanceType, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "ListInstanceTypes")
	defer span.End()

	input := &ec2.DescribeInstanceTypesInput{MaxResults: ptr.ToInt32(100)}
	pag := ec2.NewDescribeInstanceTypesPaginator(c.ec2, input)

	res := make([]types.InstanceTypeInfo, 0, 128)
	for pag.HasMorePages() {
		resp, err := pag.NextPage(ctx)
		if err != nil {
			if isAWSUnauthorizedError(err) {
				err = clients.UnauthorizedErr
			}
			span.SetStatus(codes.Error, err.Error())
			return nil, fmt.Errorf("cannot list instance types: %w", err)
		}
		res = append(res, resp.InstanceTypes...)
	}

	// convert to the client type
	instances, err := NewInstanceTypes(ctx, res)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return nil, fmt.Errorf("cannot convert instance types: %w", err)
	}

	return instances, nil
}

func (c *ec2Client) ListLaunchTemplates(ctx context.Context) ([]*clients.LaunchTemplate, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "ListLaunchTemplates")
	defer span.End()

	input := &ec2.DescribeLaunchTemplatesInput{MaxResults: ptr.ToInt32(100)}
	pag := ec2.NewDescribeLaunchTemplatesPaginator(c.ec2, input)

	var res []*clients.LaunchTemplate
	for pag.HasMorePages() {
		resp, err := pag.NextPage(ctx)
		if err != nil {
			if isAWSUnauthorizedError(err) {
				err = clients.UnauthorizedErr
			}
			span.SetStatus(codes.Error, err.Error())
			return nil, fmt.Errorf("cannot list launch templates: %w", err)
		}

		for _, awsTemplate := range resp.LaunchTemplates {
			t := clients.LaunchTemplate{
				ID:   ptr.From(awsTemplate.LaunchTemplateId),
				Name: ptr.From(awsTemplate.LaunchTemplateName),
			}
			res = append(res, &t)
		}
	}

	return res, nil
}

func (c *ec2Client) RunInstances(ctx context.Context, launchTemplateID string, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "RunInstances")
	defer span.End()

	if !c.assumed {
		return nil, nil, http.ServiceAccountUnsupportedOperationErr
	}
	logger := logger(ctx)
	logger.Trace().Msg("Run AWS EC2 instance")

	var templateSpec *types.LaunchTemplateSpecification
	if launchTemplateID != "" {
		templateSpec = &types.LaunchTemplateSpecification{
			LaunchTemplateId: ptr.To(launchTemplateID),
		}
	}

	encodedUserData := base64.StdEncoding.EncodeToString(userData)
	input := &ec2.RunInstancesInput{
		LaunchTemplate: templateSpec,
		MaxCount:       ptr.To(amount),
		MinCount:       ptr.To(amount),
		InstanceType:   instanceType,
		ImageId:        ptr.To(AMI),
		KeyName:        &keyName,
		UserData:       &encodedUserData,
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
		span.SetStatus(codes.Error, err.Error())
		return nil, nil, fmt.Errorf("cannot run instances: %w", err)
	}

	instances := c.parseRunInstancesResponse(resp)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
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

func (c *ec2Client) GetAccountId(ctx context.Context) (string, error) {
	ctx, span := otel.Tracer(TraceName).Start(ctx, "GetAccountId")
	defer span.End()

	input := &sts.GetCallerIdentityInput{}
	out, err := c.sts.GetCallerIdentity(ctx, input)
	if err != nil {
		return "", fmt.Errorf("cannot get caller's identity: %w", err)
	}

	return *out.Account, nil
}

func getRoleName(arn string) (string, error) {
	arnParts := strings.Split(arn, ":")
	if len(arnParts) == 0 {
		return "", fmt.Errorf("%w: ARN has no colons: %s", http.ARNParsingError, arn)
	}
	roleName := strings.Split(arnParts[len(arnParts)-1], "/")
	if len(roleName) != 2 {
		return "", fmt.Errorf("%w: ARN has incorrect syntax: %s rolename parsing result: %s",
			http.ARNParsingError, arn, roleName)
	} else if roleName[0] != "role" {
		return "", fmt.Errorf("%w: ARN does not have any role: %s", http.ARNParsingError, arn)
	}
	return roleName[1], nil
}

func (c *ec2Client) listAttachedRolePolicies(ctx context.Context, auth *clients.Authentication) ([]*iamTypes.AttachedPolicy, error) {
	logger := logger(ctx)
	logger.Info().Msgf("Parsing ARN to get a friendly role name.")
	roleName, err := getRoleName(auth.Payload)
	if err != nil {
		return nil, fmt.Errorf("unable to parse ARN: %w", err)
	}

	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}
	output, err := c.iam.ListAttachedRolePolicies(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("cannot list attached role policies: %w", err)
	}
	result := make([]*iamTypes.AttachedPolicy, len(output.AttachedPolicies))
	for i := range output.AttachedPolicies {
		result[i] = &output.AttachedPolicies[i]
	}
	return result, nil
}

func (c *ec2Client) CheckPermission(ctx context.Context, auth *clients.Authentication) ([]string, error) {
	logger := logger(ctx)
	logger.Info().Msgf("Listing policies attached to the role.")
	_, err := c.listAttachedRolePolicies(ctx, auth)
	if err != nil {
		return nil, err
	}

	// TODO: GetPolicies for listed policies.

	// TODO: GetPolicyVersion for policies.

	return nil, nil
}
