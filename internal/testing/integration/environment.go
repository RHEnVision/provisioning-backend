package integration

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/logging"
	"github.com/RHEnVision/provisioning-backend/internal/migrations"
	"github.com/RHEnVision/provisioning-backend/internal/queue/jq"
	"github.com/rs/zerolog/log"
)

func InitEnvironment(ctx context.Context, envPath string) {
	config.Initialize("config/test.env", envPath)
	logging.InitializeStdout()
	ctx = ctxval.WithLogger(ctx, &log.Logger)
	err := jq.Initialize(ctx, &log.Logger)
	if err != nil {
		panic(fmt.Errorf("cannot initialize job queue: %w", err))
	}
	jq.RegisterJobs(&log.Logger)
	jq.StartDequeueLoop(ctx)

	err = db.Initialize(context.Background(), "integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %w (integration schema)", err))
	}
}

func CloseEnvironment(ctx context.Context) {
	db.Close()
	jq.StopDequeueLoop(ctx)
}

func DbDrop() {
	err := migrations.Seed(context.Background(), "drop_integration")
	if err != nil {
		panic(err)
	}
}

func DbMigrate() {
	err := migrations.Migrate(context.Background(), "integration")
	if err != nil {
		panic(err)
	}
}

func DbSeed() {
	err := migrations.Seed(context.Background(), "integration")
	if err != nil {
		panic(err)
	}
}
