package ec2

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	cfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	stsTypes "github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type EC2Client struct {
	ec2     *ec2.Client
	context context.Context
	log     zerolog.Logger
}

func init() {
	clients.GetEC2Client = NewEC2Client
	clients.GetEC2ClientWithRegion = NewEC2ClientWithRegion
}

func NewEC2Client(ctx context.Context) (clients.EC2, error) {
	c, _ := NewEC2ClientWithRegion(ctx, config.AWS.Region)
	return c, nil
}

func NewEC2ClientWithRegion(ctx context.Context, region string) (clients.EC2, error) {
	c := &EC2Client{
		context: ctx,
		log:     ctxval.Logger(ctx).With().Str("client", "ec2").Logger(),
	}

	log.Trace().Msg("Creating new EC2 client")
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		config.AWS.Key, config.AWS.Secret, config.AWS.Session))

	c.ec2 = ec2.New(ec2.Options{
		Region:      region,
		Credentials: cache,
	})

	return c, nil
}

// ImportPubkey imports a key and returns AWS ID
func (c *EC2Client) ImportPubkey(key *models.Pubkey, tag string) (string, error) {
	log.Trace().Msgf("Importing AWS key-pair named '%s' with tag '%s'", key.Name, tag)
	input := &ec2.ImportKeyPairInput{}
	input.KeyName = aws.String(key.Name)
	input.PublicKeyMaterial = []byte(key.Body)
	input.TagSpecifications = []types.TagSpecification{
		{
			ResourceType: types.ResourceTypeKeyPair,
			Tags: []types.Tag{
				{
					Key:   aws.String("rhhc:id"),
					Value: aws.String(tag),
				},
			},
		},
	}
	output, err := c.ec2.ImportKeyPair(c.context, input)

	if err != nil {
		if IsOperationError(err, "InvalidKeyPair.Duplicate") {
			return "", fmt.Errorf("cannot import SSH key %s: %w", key.Name, DuplicatePubkeyErr)
		} else {
			return "", fmt.Errorf("cannot import SSH key %s: %w", key.Name, err)
		}
	}

	return aws.ToString(output.KeyPairId), nil
}

func (c *EC2Client) DeleteSSHKey(handle string) error {
	log.Trace().Msgf("Deleting AWS key-pair with handle %s", handle)
	input := &ec2.DeleteKeyPairInput{}
	input.KeyPairId = aws.String(handle)
	_, err := c.ec2.DeleteKeyPair(c.context, input)

	if err != nil {
		return fmt.Errorf("cannot delete SSH key %v: %w", input.KeyPairId, err)
	}

	return nil
}

func (c *EC2Client) CreateEC2ClientFromConfig(crd *stsTypes.Credentials) (clients.EC2, error) {
	newCfg, err := cfg.LoadDefaultConfig(c.context, cfg.WithRegion(config.AWS.Region), cfg.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(*crd.AccessKeyId, *crd.SecretAccessKey, *crd.SessionToken),
	))

	if err != nil {
		return nil, fmt.Errorf("cannot create a new ec2 config: %w", err)
	}

	newClient := &EC2Client{
		ec2:     ec2.NewFromConfig(newCfg),
		context: c.context,
		log:     ctxval.Logger(c.context).With().Str("client", "ec2").Logger(),
	}

	return newClient, nil
}

func (c *EC2Client) ListInstanceTypesWithPaginator() ([]types.InstanceTypeInfo, error) {
	input := &ec2.DescribeInstanceTypesInput{MaxResults: aws.Int32(100)}
	pag := ec2.NewDescribeInstanceTypesPaginator(c.ec2, input)

	res := make([]types.InstanceTypeInfo, 0, 128)
	for pag.HasMorePages() {
		resp, err := pag.NextPage(c.context)
		if err != nil {
			return nil, fmt.Errorf("cannot list instance types: %w", err)
		}
		res = append(res, resp.InstanceTypes...)
	}
	return res, nil
}

func (c *EC2Client) RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error) {
	log.Trace().Msg("Run AWS EC2 instance")
	encodedUserData := base64.StdEncoding.EncodeToString(userData)
	input := &ec2.RunInstancesInput{
		MaxCount:     aws.Int32(amount),
		MinCount:     aws.Int32(amount),
		InstanceType: instanceType,
		ImageId:      aws.String(AMI),
		KeyName:      &keyName,
		UserData:     &encodedUserData,
	}
	if name != nil {
		input.TagSpecifications = []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String(*name),
					},
				},
			},
		}
	}
	resp, err := c.ec2.RunInstances(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot run instances: %w", err)
	}
	instances := c.parseRunInstancesResponse(resp)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot ParseRunInstancesResponse: %w", err)
	}
	return instances, resp.ReservationId, nil
}

func (c *EC2Client) parseRunInstancesResponse(respAWS *ec2.RunInstancesOutput) []*string {
	instances := respAWS.Instances
	list := make([]*string, 0, len(instances))
	for _, instance := range instances {
		list = append(list, instance.InstanceId)
	}
	return list
}
