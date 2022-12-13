package jobs_test

import (
	"context"
	"testing"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	clientStubs "github.com/RHEnVision/provisioning-backend/internal/clients/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	daoStubs "github.com/RHEnVision/provisioning-backend/internal/dao/stubs"
	"github.com/RHEnVision/provisioning-backend/internal/jobs"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareContext(t *testing.T) context.Context {
	t.Helper()

	ctx := daoStubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = clientStubs.WithEC2Client(ctx)
	ctx = daoStubs.WithReservationDao(ctx)
	ctx = daoStubs.WithPubkeyDao(ctx)

	return ctx
}

func prepareReservation(t *testing.T, ctx context.Context, pk *models.Pubkey) *models.AWSReservation {
	t.Helper()
	detail := &models.AWSDetail{
		Region:       "us-east-1",
		InstanceType: "t1.micro",
		Amount:       1,
		PowerOff:     false,
	}
	reservation := &models.AWSReservation{
		PubkeyID: pk.ID,
		SourceID: "irrelevant",
		ImageID:  "irrelevant",
		Detail:   detail,
	}
	reservation.AccountID = 1
	reservation.Status = "Created"
	reservation.Provider = models.ProviderTypeAWS
	reservation.Steps = 2
	return reservation
}

func TestHandleEnsurePubkeyOnAWS(t *testing.T) {
	keys := []struct {
		KeyType string
		Pubkey  *models.Pubkey
	}{
		{
			"rsa",
			factories.NewPubkeyRSA(),
		},
		{
			"ed25519",
			factories.NewPubkeyED25519(),
		},
	}

	for _, testKey := range keys {
		t.Run("exists_"+testKey.KeyType, func(t *testing.T) {
			ctx := prepareContext(t)
			pkt := factories.PubkeyWithTrans(t, ctx, testKey.Pubkey)

			// using a static key for which we know it's real AWS fingerprint
			pk := &models.Pubkey{
				Name: "provisioningName",
				Body: pkt.Body,
			}
			err := daoStubs.AddPubkey(ctx, pk)
			require.NoError(t, err, "failed to add stubbed key")

			err = clientStubs.AddStubbedEC2KeyPair(ctx, &types.KeyPairInfo{
				KeyName:        ptr.To("awsName"),
				KeyFingerprint: ptr.To(pkt.FindAwsFingerprint(ctx)),
				KeyType:        types.KeyType(testKey.KeyType),
				PublicKey:      &pk.Body,
				Tags: []types.Tag{
					{Key: ptr.To("rhhc:id"), Value: ptr.To("awsPreviouslyImportedTag")},
				},
			})
			require.NoError(t, err, "failed to add stubbed key to ec2 stub")

			reservation := prepareReservation(t, ctx, pk)
			rDao := dao.GetReservationDao(ctx)
			err = rDao.CreateAWS(ctx, reservation)
			require.NoError(t, err, "failed to add stubbed reservation")

			args := &jobs.LaunchInstanceAWSTaskArgs{
				ReservationID: reservation.ID,
				Region:        reservation.Detail.Region,
				PubkeyID:      pk.ID,
				SourceID:      reservation.SourceID,
				Detail:        reservation.Detail,
				ARN:           &clients.Authentication{ProviderType: models.ProviderTypeAWS, Payload: "arn:aws:123123123123"},
			}

			err = jobs.DoEnsurePubkeyOnAWS(ctx, args)
			require.NoError(t, err, "the ensure pubkey job failed to run")

			resAfter, err := rDao.GetAWSById(ctx, reservation.ID)
			require.NoError(t, err, "failed to add stubbed reservation")
			assert.Equal(t, "awsName", resAfter.Detail.PubkeyName)

			pkDao := dao.GetPubkeyDao(ctx)
			pkrList, err := pkDao.UnscopedListResourcesByPubkeyId(ctx, pk.ID)
			require.NoError(t, err)
			assert.Equal(t, 1, len(pkrList))
		})
	}

	t.Run("not_exists", func(t *testing.T) {
		ctx := prepareContext(t)

		pk := &models.Pubkey{
			Name: factories.SeqNameWithPrefix("pubkey"),
			Body: factories.GenerateRSAPubKey(t),
		}
		err := daoStubs.AddPubkey(ctx, pk)
		require.NoError(t, err, "failed to add stubbed key")

		reservation := prepareReservation(t, ctx, pk)
		rDao := dao.GetReservationDao(ctx)
		err = rDao.CreateAWS(ctx, reservation)
		require.NoError(t, err, "failed to add stubbed reservation")

		args := &jobs.LaunchInstanceAWSTaskArgs{
			ReservationID: reservation.ID,
			Region:        reservation.Detail.Region,
			PubkeyID:      pk.ID,
			SourceID:      reservation.SourceID,
			Detail:        reservation.Detail,
			ARN:           &clients.Authentication{ProviderType: models.ProviderTypeAWS, Payload: "arn:aws:123123123123"},
		}

		err = jobs.DoEnsurePubkeyOnAWS(ctx, args)
		require.NoError(t, err, "the ensure pubkey job failed to run")

		resAfter, err := rDao.GetAWSById(ctx, reservation.ID)
		require.NoError(t, err)
		assert.Equal(t, pk.Name, resAfter.Detail.PubkeyName)

		pkDao := dao.GetPubkeyDao(ctx)
		pkrList, err := pkDao.UnscopedListResourcesByPubkeyId(ctx, pk.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pkrList))
	})

	t.Run("Pubkey on AWS and resource exists - should copy the Handle and Tag", func(t *testing.T) {
		ctx := prepareContext(t)

		pk := &models.Pubkey{
			Name: factories.SeqNameWithPrefix("pubkey"),
			Body: factories.GenerateRSAPubKey(t),
		}
		pk = factories.PubkeyWithTrans(t, ctx, pk)
		err := daoStubs.AddPubkey(ctx, pk)
		require.NoError(t, err, "failed to add stubbed key")

		pkDao := dao.GetPubkeyDao(ctx)
		err = pkDao.UnscopedCreateResource(ctx, &models.PubkeyResource{
			PubkeyID: pk.ID,
			Tag:      "veryRandom",
			Provider: models.ProviderTypeAWS,
			Handle:   "someAwsHandle",
			Region:   "us-east-1",
		})
		require.NoError(t, err, "failed to create stubbed pubkey resource")

		err = clientStubs.AddStubbedEC2KeyPair(ctx, &types.KeyPairInfo{
			KeyPairId:      ptr.To("veryCoolAwsId"),
			KeyName:        ptr.To("awsName"),
			KeyFingerprint: ptr.To(pk.FindAwsFingerprint(ctx)),
			KeyType:        "rsa",
			PublicKey:      &pk.Body,
			Tags: []types.Tag{
				{Key: ptr.To("rhhc:id"), Value: ptr.To("awsPreviouslyImportedTag")},
			},
		})
		require.NoError(t, err, "failed to add stubbed key to ec2 stub")

		reservation := prepareReservation(t, ctx, pk)
		rDao := dao.GetReservationDao(ctx)
		err = rDao.CreateAWS(ctx, reservation)
		require.NoError(t, err, "failed to add stubbed reservation")

		args := &jobs.LaunchInstanceAWSTaskArgs{
			ReservationID: reservation.ID,
			Region:        reservation.Detail.Region,
			PubkeyID:      pk.ID,
			SourceID:      reservation.SourceID,
			Detail:        reservation.Detail,
			ARN:           &clients.Authentication{ProviderType: models.ProviderTypeAWS, Payload: "arn:aws:123123123123"},
		}

		err = jobs.DoEnsurePubkeyOnAWS(ctx, args)
		require.NoError(t, err, "the ensure pubkey job failed to run")

		pkrList, err := pkDao.UnscopedListResourcesByPubkeyId(ctx, pk.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pkrList))
		assert.Equal(t, "awsPreviouslyImportedTag", pkrList[0].Tag)
		assert.Equal(t, "veryCoolAwsId", pkrList[0].Handle)
	})
}
