package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func ImportSSHKey(ctx context.Context, body string) (string, error) {
	input := &ec2.ImportKeyPairInput{}
	input.KeyName = aws.String("Red Hat Portal Key")
	input.PublicKeyMaterial = []byte(body)
	output, err := EC2.ImportKeyPair(ctx, input)

	if err != nil {
		return "", err
	}

	return aws.ToString(output.KeyPairId), nil
}

func DeleteSSHKey(ctx context.Context, cid string) error {
	input := &ec2.DeleteKeyPairInput{}
	input.KeyPairId = aws.String(cid)
	_, err := EC2.DeleteKeyPair(ctx, input)

	if err != nil {
		return err
	}

	return nil
}
