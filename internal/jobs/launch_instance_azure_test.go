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

func prepareAzureContext(t *testing.T) context.Context {
	t.Helper()

	ctx := daoStubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = clientStubs.WithAzureClient(ctx)
	ctx = daoStubs.WithReservationDao(ctx)
	ctx = daoStubs.WithPubkeyDao(ctx)

	return ctx
}

func prepareAzureReservation(t *testing.T, ctx context.Context, pk *models.Pubkey) *models.AzureReservation {
	t.Helper()

	detail := &models.AzureDetail{
		Location:     "useast",
		InstanceSize: "Basic_A0",
		Amount:       1,
		PowerOff:     false,
	}
	reservation := &models.AzureReservation{
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

func TestDoEnsureAzureResourceGroup(t *testing.T) {
	ctx := prepareAzureContext(t)

	pk := factories.NewPubkeyRSA()
	err := daoStubs.AddPubkey(ctx, pk)
	require.NoError(t, err, "failed to add stubbed key")

	res := prepareAzureReservation(t, ctx, pk)

	rDao := dao.GetReservationDao(ctx)
	err = rDao.CreateAzure(ctx, res)
	require.NoError(t, err, "failed to add stubbed reservation")

	args := &jobs.LaunchInstanceAzureTaskArgs{
		AzureImageID:  "/subscriptions/subUUID/rgName/images/uuid2",
		Location:      "useast",
		PubkeyID:      pk.ID,
		ReservationID: res.ID,
		SourceID:      "2",
		Subscription:  clients.NewAuthentication("subUUID", models.ProviderTypeAzure),
	}

	err = jobs.DoEnsureAzureResourceGroup(ctx, args)
	require.NoError(t, err, "the ensure resource group failed to run")

	assert.True(t, clientStubs.DidCreateAzureResourceGroup(ctx, "redhat-deployed"))
}

func TestDoLaunchInstanceAzure(t *testing.T) {
	ctx := prepareAzureContext(t)

	pk := factories.NewPubkeyRSA()
	err := daoStubs.AddPubkey(ctx, pk)
	require.NoError(t, err, "failed to add stubbed key")

	res := prepareAzureReservation(t, ctx, pk)
	res.Detail.Amount = 2

	rDao := dao.GetReservationDao(ctx)
	err = rDao.CreateAzure(ctx, res)
	require.NoError(t, err, "failed to add stubbed reservation")

	args := &jobs.LaunchInstanceAzureTaskArgs{
		AzureImageID:  "/subscriptions/subUUID/rgName/images/uuid2",
		Location:      "useast",
		PubkeyID:      pk.ID,
		ReservationID: res.ID,
		SourceID:      "2",
		Subscription:  clients.NewAuthentication("subUUID", models.ProviderTypeAzure),
	}

	err = jobs.DoLaunchInstanceAzure(ctx, args)
	require.NoError(t, err, "launch instances failed to run")

	assert.Equal(t, 2, clientStubs.CountStubAzureVMs(ctx))
}
