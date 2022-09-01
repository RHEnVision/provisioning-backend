package sts

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/aws/aws-sdk-go-v2/aws"
	con "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/rs/zerolog"
)

type STSClient struct {
	sts *sts.Client
	ctx context.Context
	log zerolog.Logger
}

func init() {
	clients.GetSTSClient = NewSTSClient
}

func NewSTSClient(ctx context.Context) (clients.STS, error) {
	c := STSClient{
		ctx: ctx,
		log: ctxval.Logger(ctx).With().Str("client", "sts").Logger(),
	}

	cfg, err := con.LoadDefaultConfig(ctx, con.WithRegion(config.AWS.Region),
		con.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(config.AWS.Key, config.AWS.Secret, config.AWS.Session)))
	if err != nil {
		c.log.Error().Err(err).Msgf("cannot create an sts client %s", err)
		return nil, fmt.Errorf("cannot create an sts client %w", err)
	}
	c.sts = sts.NewFromConfig(cfg)
	return &c, nil
}

func (c *STSClient) AssumeRole(arn string) (*types.Credentials, error) {
	output, err := c.sts.AssumeRole(c.ctx, &sts.AssumeRoleInput{
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String("name"),
	})

	if err != nil {
		c.log.Error().Err(err).Msgf("cannot assume role %s", err)
		return nil, fmt.Errorf("cannot assume role %w", err)
	}

	return output.Credentials, nil
}
