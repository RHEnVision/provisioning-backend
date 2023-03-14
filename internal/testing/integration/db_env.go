package integration

import (
	"context"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/db"
	"github.com/RHEnVision/provisioning-backend/internal/migrations"
)

func InitDbEnvironment(ctx context.Context) {
	err := db.Initialize(context.Background(), "integration")
	if err != nil {
		panic(fmt.Errorf("cannot connect to database: %w (integration schema)", err))
	}
}

func CloseDbEnvironment(ctx context.Context) {
	db.Close()
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
