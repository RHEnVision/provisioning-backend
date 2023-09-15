package stubs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/clients/http"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ec2CtxKeyType int

const ec2CtxKey ec2CtxKeyType = iota

type EC2ClientStub struct {
	Imported []*types.KeyPairInfo
}

func init() {
	clients.GetEC2Client = newEC2CustomerClientStubWithRegion
	clients.GetServiceEC2Client = newEC2ServiceClientStubWithRegion
}

func WithEC2Client(parent context.Context) context.Context {
	ctx := context.WithValue(parent, ec2CtxKey, &EC2ClientStub{})
	return ctx
}

func AddStubbedEC2KeyPair(ctx context.Context, info *types.KeyPairInfo) error {
	si, err := getEC2StubFromContext(ctx)
	if err != nil {
		return err
	}
	si.Imported = append(si.Imported, info)
	return nil
}

func newEC2ServiceClientStubWithRegion(ctx context.Context, region string) (clients.EC2, error) {
	return nil, nil
}

func newEC2CustomerClientStubWithRegion(ctx context.Context, _ *clients.Authentication, _ string) (si clients.EC2, err error) {
	return getEC2StubFromContext(ctx)
}

func getEC2StubFromContext(ctx context.Context) (*EC2ClientStub, error) {
	var si *EC2ClientStub
	var err error
	var ok bool
	if si, ok = ctx.Value(ec2CtxKey).(*EC2ClientStub); !ok {
		err = ContextReadError
	}
	return si, err
}

func (mock *EC2ClientStub) Status(ctx context.Context) error {
	return nil
}

func (mock *EC2ClientStub) ImportPubkey(ctx context.Context, key *models.Pubkey, tag string) (string, error) {
	ec2KeyID := fmt.Sprintf("key-%d", len(mock.Imported))
	fingerprint := key.FindAwsFingerprint(ctx)
	keyName := key.Name // copy the name
	ec2Key := &types.KeyPairInfo{
		KeyName: &keyName,

		KeyFingerprint: &fingerprint,
		KeyPairId:      &ec2KeyID,
		PublicKey:      &key.Body,
		Tags: []types.Tag{{
			Key:   ptr.To("rh-kid"),
			Value: &tag,
		}},
	}
	mock.Imported = append(mock.Imported, ec2Key)
	return *ec2Key.KeyPairId, nil
}

func (mock *EC2ClientStub) GetPubkeyName(ctx context.Context, fingerprint string) (string, error) {
	for _, key := range mock.Imported {
		if *key.KeyFingerprint == fingerprint {
			return *key.KeyName, nil
		}
	}
	return "", http.PubkeyNotFoundErr
}

func (mock *EC2ClientStub) DeleteSSHKey(ctx context.Context, handle string) error {
	return nil
}

func (mock *EC2ClientStub) ListAllRegions(ctx context.Context) ([]clients.Region, error) {
	return []clients.Region{
		"us-east-1",
		"eu-central-1",
	}, nil
}

func (mock *EC2ClientStub) ListAllZones(ctx context.Context, region clients.Region) ([]clients.Zone, error) {
	return []clients.Zone{
		"us-east-1a",
		"us-east-1b",
		"us-east-1c",
		"eu-central-1a",
		"eu-central-1b",
		"eu-central-1c",
	}, nil
}

func (mock *EC2ClientStub) ListInstanceTypes(ctx context.Context) ([]*clients.InstanceType, error) {
	return []*clients.InstanceType{
		{
			Name:               "t4g.nano",
			VCPUs:              2,
			Cores:              2,
			MemoryMiB:          500,
			EphemeralStorageGB: 0,
			Supported:          false,
			Architecture:       clients.ArchitectureTypeArm64,
		},
		{
			Name:               "a1.2xlarge",
			VCPUs:              8,
			Cores:              8,
			MemoryMiB:          16000,
			EphemeralStorageGB: 0,
			Supported:          true,
			Architecture:       clients.ArchitectureTypeX86_64,
		},
		{
			Name:               "c5.xlarge",
			VCPUs:              4,
			Cores:              4,
			MemoryMiB:          8000,
			EphemeralStorageGB: 0,
			Supported:          true,
			Architecture:       clients.ArchitectureTypeX86_64,
		},
	}, nil
}

func (mock *EC2ClientStub) ListLaunchTemplates(ctx context.Context) ([]*clients.LaunchTemplate, error) {
	return []*clients.LaunchTemplate{
		{
			ID:   "lt-8732678436272377",
			Name: "Nano ARM64 load balancer",
		},
		{
			ID:   "lt-8732678438462378",
			Name: "XXLarge AMD64 database",
		},
	}, nil
}

func (mock *EC2ClientStub) CheckPermission(ctx context.Context, auth *clients.Authentication) ([]string, error) {
	return nil, nil
}

func (mock *EC2ClientStub) RunInstances(ctx context.Context, details *clients.AWSInstanceParams, amount int32, name string, reservation *models.AWSReservation) ([]*string, *string, error) {
	return nil, nil, nil
}

func (mock *EC2ClientStub) GetAccountId(ctx context.Context) (string, error) {
	return "", nil
}

func (mock *EC2ClientStub) DescribeInstanceDetails(ctx context.Context, InstanceIds []string) ([]*clients.InstanceDescription, error) {
	id := "i-0a4caa2cf5b097ce1"
	dns := "ec2-51-83-81-17.compute-1.amazonaws.com"
	ip := "54.11.88.17"
	return []*clients.InstanceDescription{
		{
			ID:         id,
			PublicDNS:  dns,
			PublicIPv4: ip,
		},
	}, nil
}
