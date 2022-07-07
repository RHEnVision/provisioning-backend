package ec2

import (
	"context"
	"fmt"

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

type Client struct {
	ec2     *ec2.Client
	context context.Context
	log     zerolog.Logger
}

func NewEC2Client(ctx context.Context) *Client {
	c := Client{
		context: ctx,
		log:     ctxval.GetLogger(ctx).With().Str("client", "ec2").Logger(),
	}

	log.Trace().Msg("Creating new EC2 client")
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		config.AWS.Key, config.AWS.Secret, config.AWS.Session))

	c.ec2 = ec2.New(ec2.Options{
		Region:      config.AWS.Region,
		Credentials: cache,
	})

	return &c
}

// ImportPubkey imports a key and returns AWS ID
func (c *Client) ImportPubkey(key *models.Pubkey, tag string) (string, error) {
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

func (c *Client) DeleteSSHKey(handle string) error {
	log.Trace().Msgf("Deleting AWS key-pair with handle %s", handle)
	input := &ec2.DeleteKeyPairInput{}
	input.KeyPairId = aws.String(handle)
	_, err := c.ec2.DeleteKeyPair(c.context, input)

	if err != nil {
		return fmt.Errorf("cannot delete SSH key %v: %w", input.KeyPairId, err)
	}

	return nil
}

func (c *Client) CreateEC2ClientFromConfig(crd *stsTypes.Credentials) (*Client, error) {
	newCfg, err := cfg.LoadDefaultConfig(c.context, cfg.WithRegion(config.AWS.Region), cfg.WithCredentialsProvider(
		credentials.NewStaticCredentialsProvider(*crd.AccessKeyId, *crd.SecretAccessKey, *crd.SessionToken),
	))

	if err != nil {
		return nil, fmt.Errorf("cannot create a new ec2 config: %w", err)
	}

	newClient := &Client{
		ec2:     ec2.NewFromConfig(newCfg),
		context: c.context,
		log:     ctxval.GetLogger(c.context).With().Str("client", "ec2").Logger(),
	}

	return newClient, nil
}

func (c *Client) ListInstanceTypes() ([]types.InstanceTypeInfo, error) {
	log.Trace().Msg("Listing AWS EC2 instance types")
	input := &ec2.DescribeInstanceTypesInput{
		MaxResults: aws.Int32(100),
	}

	resp, err := c.ec2.DescribeInstanceTypes(c.context, input)
	if err != nil {
		return nil, fmt.Errorf("cannot list instance types: %w", err)
	}

	log.Debug().Msgf("Number AWS EC2 instance types: %d", len(resp.InstanceTypes))
	if len(resp.InstanceTypes) == 100 {
		return nil, OperationNotPermittedErr
	}

	return resp.InstanceTypes, nil
}
