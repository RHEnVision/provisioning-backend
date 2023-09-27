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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareGCPContext(t *testing.T) context.Context {
	t.Helper()
	ctx := daoStubs.WithAccountDaoOne(context.Background())
	ctx = identity.WithTenant(t, ctx)
	ctx = clientStubs.WithGCPCCustomerClient(ctx)
	ctx = daoStubs.WithReservationDao(ctx)
	ctx = daoStubs.WithPubkeyDao(ctx)
	return ctx
}

func prepareGCPReservation(t *testing.T, ctx context.Context, pk *models.Pubkey) *models.GCPReservation {
	t.Helper()

	detail := &models.GCPDetail{
		Zone:        "europe-west8-c",
		MachineType: "e2-micro",
		NamePattern: ptr.To("instance-#####"),
		Amount:      1,
		PowerOff:    false,
	}
	reservation := &models.GCPReservation{
		PubkeyID: pk.ID,
		SourceID: "irrelevant",
		ImageID:  "irrelevant",
		Detail:   detail,
	}
	reservation.AccountID = 1
	reservation.Status = "Created"
	reservation.Provider = models.ProviderTypeGCP
	reservation.Steps = 2
	return reservation
}

func TestDoLaunchInstanceGCP(t *testing.T) {
	ctx := prepareGCPContext(t)

	pk := factories.NewPubkeyRSA()
	err := daoStubs.AddPubkey(ctx, pk)
	require.NoError(t, err, "failed to add stubbed key")

	res := prepareGCPReservation(t, ctx, pk)
	res.Detail.Amount = 2
	rDao := dao.GetReservationDao(ctx)
	err = rDao.CreateGCP(ctx, res)
	require.NoError(t, err, "failed to add stubbed reservation")

	args := &jobs.LaunchInstanceGCPTaskArgs{
		ImageName:     "composer-api-3b6225fc-d55a-4dcc-9d0a-b478ae152a",
		Zone:          "europe-west8-c",
		PubkeyID:      pk.ID,
		ReservationID: res.ID,
		ProjectID:     clients.NewAuthentication("example-project-id", models.ProviderTypeGCP),
		Detail:        res.Detail,
	}

	err = jobs.DoLaunchInstanceGCP(ctx, args)
	require.NoError(t, err, "launch instances failed to run")
	assert.Equal(t, 2, clientStubs.CountStubInstancesGCP(ctx))

	t.Run("fetch instances description", func(t *testing.T) {
		err = jobs.FetchInstancesDescriptionGCP(ctx, args)
		require.NoError(t, err, "fetch instances description failed to run")

		resultInstances, err := rDao.ListInstances(ctx, res.ID)
		require.NoError(t, err, "failed to fetch created instances")
		assert.Equal(t, 2, len(resultInstances))
		assert.NotEmpty(t, resultInstances[0].Detail.PublicIPv4)
		assert.Equal(t, "10.0.0.10", resultInstances[0].Detail.PublicIPv4)
		assert.Equal(t, "10.0.0.11", resultInstances[1].Detail.PublicIPv4)
	})
}
