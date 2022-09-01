package cloudwatchlogs

import (
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// CWL is Cloudwatch AWS client
var CWL *cloudwatchlogs.Client

func Initialize() {
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		config.Cloudwatch.Key, config.Cloudwatch.Secret, config.Cloudwatch.Session))

	CWL = cloudwatchlogs.New(cloudwatchlogs.Options{
		Region:      config.Cloudwatch.Region,
		Credentials: cache,
	})
}
