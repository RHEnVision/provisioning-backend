package aws

import (
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2 client will be dropped once we implement Sources API
var EC2 *ec2.Client

// CWL is Cloudwatch AWS client
var CWL *cloudwatchlogs.Client

func Initialize() {
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		config.Cloudwatch.Key, config.Cloudwatch.Secret, config.Cloudwatch.Session))

	EC2 = ec2.New(ec2.Options{
		Region:      config.Cloudwatch.Region,
		Credentials: cache,
	})
	CWL = cloudwatchlogs.New(cloudwatchlogs.Options{
		Region:      config.Cloudwatch.Region,
		Credentials: cache,
	})
}
