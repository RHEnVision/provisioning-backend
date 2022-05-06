package aws

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

var EC2 *ec2.Client
var CWL *cloudwatchlogs.Client

func Initialize() {
	key := os.Getenv("AWS_KEY")
	secret := os.Getenv("AWS_SECRET")
	session := os.Getenv("AWS_SESSION")
	region := os.Getenv("AWS_REGION")
	cache := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(key, secret, session))

	EC2 = ec2.New(ec2.Options{
		Region:      region,
		Credentials: cache,
	})
	CWL = cloudwatchlogs.New(cloudwatchlogs.Options{
		Region:      region,
		Credentials: cache,
	})
}
