package ec2

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
	log.Trace().Msgf("Importing pubkey '%s'", key.Name)
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
		return "", fmt.Errorf("cannot import SSH key %s: %w", key.Name, err)
	}

	return aws.ToString(output.KeyPairId), nil
}

func (c *Client) DeleteSSHKey(cid string) error {
	log.Trace().Msgf("Deleting pubkey with cid %s", cid)
	input := &ec2.DeleteKeyPairInput{}
	input.KeyPairId = aws.String(cid)
	_, err := c.ec2.DeleteKeyPair(c.context, input)

	if err != nil {
		return fmt.Errorf("cannot delete SSH key %v: %w", input.KeyPairId, err)
	}

	return nil
}
