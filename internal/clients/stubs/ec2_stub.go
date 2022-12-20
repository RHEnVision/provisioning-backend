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

type ec2CtxKeyType string

var ec2CtxKey ec2CtxKeyType = "ec2-interface"

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
		err = &contextReadError{}
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
			Key:   ptr.To("rhhc:id"),
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

func (mock *EC2ClientStub) ListInstanceTypesWithPaginator(ctx context.Context) ([]*clients.InstanceType, error) {
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

func (mock *EC2ClientStub) RunInstances(ctx context.Context, name *string, amount int32, instanceType types.InstanceType, AMI string, keyName string, userData []byte) ([]*string, *string, error) {
	return nil, nil, nil
}

func (mock *EC2ClientStub) GetAccountId(ctx context.Context) (string, error) {
	return "", nil
}
