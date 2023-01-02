package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/jobs/queue"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
)

type LaunchInstanceGCPTaskArgs struct {
	// Associated reservation
	ReservationID int64 `json:"reservation_id"`

	// Associated account
	AccountID int64 `json:"account_id"`

	// Zone to provision the instances into
	Zone string `json:"zone"`

	// Associated public key
	PubkeyID int64 `json:"pubkey_id"`

	// Detail information
	Detail *models.GCPDetail `json:"detail"`

	// GCP image name as fetched from image builder
	ImageName string `json:"image_name"`

	// The project id from Sources which is linked to a specific source
	ProjectID *clients.Authentication `json:"project_id"`
}

var LaunchInstanceGCPTask = queue.RegisterTask("launch_instance_gcp", func(ctx context.Context, args LaunchInstanceGCPTaskArgs) error {
	ctx = prepareContext(ctx, "launch_instance_gcp", args, args.AccountID, args.ReservationID)
	err := HandleLaunchInstanceGCP(ctx, args)
	finishStep(ctx, args.ReservationID, err)
	return err
})

func EnqueueLaunchInstanceGCP(ctx context.Context, args LaunchInstanceGCPTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)

	t := LaunchInstanceGCPTask.WithArgs(ctx, args)
	ctxLogger.Debug().Str("tid", t.ID).Msg("Adding launch GCP job task")
	err := queue.JobQueue.Add(t)
	if err != nil {
		return fmt.Errorf("unable to enqueue task: %w", err)
	}

	return nil
}

func HandleLaunchInstanceGCP(ctx context.Context, args LaunchInstanceGCPTaskArgs) error {
	ctxLogger := ctxval.Logger(ctx)
	ctxLogger.Debug().Msg("Started launch instance GCP job")

	ctx = ctxval.WithAccountId(ctx, args.AccountID)
	logger := ctxLogger.With().Int64("reservation", args.ReservationID).Logger()
	logger.Info().Interface("args", args).Msg("Processing launch instance GCP job")

	// status updates before and after the code logic
	updateStatusBefore(ctx, args.ReservationID, "Launching instance(s)")
	defer updateStatusAfter(ctx, args.ReservationID, "Launched instance(s)", 1)

	pkD := dao.GetPubkeyDao(ctx)

	pk, err := pkD.GetById(ctx, args.PubkeyID)
	if err != nil {
		return fmt.Errorf("cannot get pubkey by id: %w", err)
	}

	gcpClient, err := clients.GetGCPClient(ctx, args.ProjectID)
	if err != nil {
		return fmt.Errorf("cannot get gcp client: %w", err)
	}

	opName, err := gcpClient.InsertInstances(ctx, ptr.To("inst-####"), &args.ImageName, args.Detail.Amount, args.Detail.MachineType, args.Zone, pk.Body)
	if err != nil {
		return fmt.Errorf("cannot run instances for gcp client: %w", err)
	}

	rDao := dao.GetReservationDao(ctx)

	err = rDao.UpdateOperationNameForGCP(ctx, args.ReservationID, *opName)
	if err != nil {
		return fmt.Errorf("cannot update operation name for GCP : %w", err)
	}

	return nil
}
