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
	"github.com/RHEnVision/provisioning-backend/internal/testing/factories"
	"github.com/RHEnVision/provisioning-backend/internal/testing/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubAWSJob struct {
	body *jobs.EnsurePubkeyOnAWSTaskArgs
}

func (s stubAWSJob) Type() string { return "stub AWS pubkey job" }
func (s stubAWSJob) Decode(out interface{}) error {
	typed := out.(*jobs.EnsurePubkeyOnAWSTaskArgs)
	typed.ARN = s.body.ARN
	typed.AccountID = s.body.AccountID
	typed.ReservationID = s.body.ReservationID
	typed.SourceID = s.body.SourceID
	typed.PubkeyID = s.body.PubkeyID
	typed.Region = s.body.Region
	return nil
}

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
	t.Run("KeyPair exists on AWS", func(t *testing.T) {
		ctx := prepareContext(t)

		pk := &models.Pubkey{
			Name: "provisioningName",
			Body: factories.GenerateRSAPubKey(t),
		}
		err := daoStubs.AddPubkey(ctx, pk)
		require.NoError(t, err, "failed to add stubbed key")

		authentication := clients.Authentication{ProviderType: models.ProviderTypeAWS, Payload: "arn:aws:123123123123"}

		ec2Client, err := clients.GetEC2Client(ctx, &authentication, "us-east-1")
		require.NoError(t, err, "failed to get stubbed EC2 client")
		pk.Name = "awsName" // change the name to get imported in the stub under
		_, err = ec2Client.ImportPubkey(ctx, pk, "some-tag")
		require.NoError(t, err, "failed to ImportPubkey to the stub")
		pk.Name = "provisioningName" // change the name back

		reservation := prepareReservation(t, ctx, pk)
		rDao := dao.GetReservationDao(ctx)
		err = rDao.CreateAWS(ctx, reservation)
		require.NoError(t, err, "failed to add stubbed reservation")

		stubbedEnsureJob := stubAWSJob{
			body: &jobs.EnsurePubkeyOnAWSTaskArgs{
				AccountID:     1,
				ReservationID: reservation.ID,
				Region:        reservation.Detail.Region,
				PubkeyID:      pk.ID,
				SourceID:      reservation.SourceID,
				ARN:           &authentication,
			},
		}

		err = jobs.HandleEnsurePubkeyOnAWS(ctx, stubbedEnsureJob)
		require.NoError(t, err, "the ensure pubkey job failed to run")

		resAfter, err := rDao.GetAWSById(ctx, reservation.ID)
		require.NoError(t, err, "failed to add stubbed reservation")
		assert.Equal(t, "awsName", resAfter.Detail.PubkeyName)

		pkDao := dao.GetPubkeyDao(ctx)
		pkrList, err := pkDao.UnscopedListResourcesByPubkeyId(ctx, pk.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, len(pkrList))
	})

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

		stubbedEnsureJob := stubAWSJob{
			body: &jobs.EnsurePubkeyOnAWSTaskArgs{
				AccountID:     1,
				ReservationID: reservation.ID,
				Region:        reservation.Detail.Region,
				PubkeyID:      pk.ID,
				SourceID:      reservation.SourceID,
				ARN:           &clients.Authentication{ProviderType: models.ProviderTypeAWS, Payload: "arn:aws:123123123123"},
			},
		}

		err = jobs.HandleEnsurePubkeyOnAWS(ctx, stubbedEnsureJob)
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
