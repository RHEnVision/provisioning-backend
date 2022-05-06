package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func ImportSSHKey(ctx context.Context, body string) (string, error) {
	input := &ec2.ImportKeyPairInput{}
	input.KeyName = aws.String("Red Hat Portal Key")
	input.PublicKeyMaterial = []byte(body)
	output, err := EC2.ImportKeyPair(ctx, input)

	if err != nil {
		return "", fmt.Errorf("cannot import SSH key %v: %w", input.PublicKeyMaterial, err)
	}

	return aws.ToString(output.KeyPairId), nil
}

func DeleteSSHKey(ctx context.Context, cid string) error {
	input := &ec2.DeleteKeyPairInput{}
	input.KeyPairId = aws.String(cid)
	_, err := EC2.DeleteKeyPair(ctx, input)

	if err != nil {
		return fmt.Errorf("cannot delete SSH key %v: %w", input.KeyPairId, err)
	}

	return nil
}
