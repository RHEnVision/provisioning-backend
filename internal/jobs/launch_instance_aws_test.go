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
		KeyType        string
		Body           string
		AwsFingerprint string
	}{
		{
			"rsa",
			"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDABkLdmd4mnZTOKcZpi0cu+YWbnVbbJSjHW5FDc4p9AgeVZA2WKx2f2x4YPOdB9NtAuVuFUpDAvBiV96Coy0747I2RXrR5abzVWbW+bIJOJXqCCBLHlUEj13CduIs40pHwVXGRdlwLZk4rChWKzg+C6sNBq5lXxtBLfKmf5S2LbhWKTfZje7OI2We2pZXiRZg58IVIA2mpvNr3MxaoMlEK92VwiVzlwOaCKbG4Ere5M1ug/5RRSXXLQjPBc6ePqg1PiVHrx2DP2jDJsGETQGj13zzqI4nvXbcu7EM/TiCJpreHZIxDgn97AtYj2IJbscz+6/aWyBlZG0be8oaLLlYN",
			"93:fb:9e:04:da:6e:5d:37:5d:2e:f4:0b:39:ef:6a:08",
		},
		{
			"ed25519",
			"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhnn80ZywmjeBFFOGm+cm+5HUwm62qTVnjKlOdYFLHN lzap",
			"gL/y6MvNmJ8jDXtsL/oMmK8jUuIefN39BBuvYw/Rndk=",
		},
	}

	for _, testKey := range keys {
		t.Run("KeyPair exists on AWS", func(t *testing.T) {
			ctx := prepareContext(t)

			// using a static key for which we know it's real AWS fingerprint
			pk := &models.Pubkey{
				Name: "provisioningName",
				Body: testKey.Body,
			}
			err := daoStubs.AddPubkey(ctx, pk)
			require.NoError(t, err, "failed to add stubbed key")

			err = clientStubs.AddStubbedEC2KeyPair(ctx, &types.KeyPairInfo{
				KeyName:        ptr.To("awsName"),
				KeyFingerprint: &testKey.AwsFingerprint,
				KeyType:        types.KeyType(testKey.KeyType),
				PublicKey:      &pk.Body,
			})
			require.NoError(t, err, "failed to add stubbed key to ec2 stub")

			reservation := prepareReservation(t, ctx, pk)
			rDao := dao.GetReservationDao(ctx)
			err = rDao.CreateAWS(ctx, reservation)
			require.NoError(t, err, "failed to add stubbed reservation")

			args := &jobs.LaunchInstanceAWSTaskArgs{
				AccountID:     1,
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

	t.Run("pubkey not on AWS", func(t *testing.T) {
		ctx := prepareContext(t)

		pk := &models.Pubkey{
			Name: factories.GetSequenceName("pubkey"),
			Body: factories.GenerateRSAPubKey(t),
		}
		err := daoStubs.AddPubkey(ctx, pk)
		require.NoError(t, err, "failed to add stubbed key")

		reservation := prepareReservation(t, ctx, pk)
		rDao := dao.GetReservationDao(ctx)
		err = rDao.CreateAWS(ctx, reservation)
		require.NoError(t, err, "failed to add stubbed reservation")

		args := &jobs.LaunchInstanceAWSTaskArgs{
			AccountID:     1,
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
}
