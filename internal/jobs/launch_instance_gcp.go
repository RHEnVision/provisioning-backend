package jobs

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	_ "github.com/RHEnVision/provisioning-backend/internal/clients/http/gcp"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/dao"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/RHEnVision/provisioning-backend/internal/ptr"
	"github.com/RHEnVision/provisioning-backend/pkg/worker"
)

type LaunchInstanceGCPTaskArgs struct {
	// Associated reservation
	ReservationID int64

	// Zone to provision the instances into
	Zone string

	// Associated public key
	PubkeyID int64

	// Detail information
	Detail *models.GCPDetail

	// GCP image name as fetched from image builder
	ImageName string

	// The project id from Sources which is linked to a specific source
	ProjectID *clients.Authentication
}

// Unmarshall arguments and handle error
func HandleLaunchInstanceGCP(ctx context.Context, job *worker.Job) {
	args, ok := job.Args.(LaunchInstanceGCPTaskArgs)
	if !ok {
		ctxval.Logger(ctx).Error().Msgf("Type assertion error for job %s, unable to finish reservation: %#v", job.ID, job.Args)
		return
	}

	logger := ctxval.Logger(ctx).With().Int64("reservation_id", args.ReservationID).Logger()
	ctx = ctxval.WithLogger(ctx, &logger)

	jobErr := DoLaunchInstanceGCP(ctx, &args)

	finishJob(ctx, args.ReservationID, jobErr)
}

// Job logic, when error is returned the job status is updated accordingly
func DoLaunchInstanceGCP(ctx context.Context, args *LaunchInstanceGCPTaskArgs) error {
	logger := ctxval.Logger(ctx)
	logger.Debug().Msg("Started launch instance GCP job")
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
