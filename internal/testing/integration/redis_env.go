package integration

import (
	"context"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
)

func InitRedisEnvironment(_ context.Context) {
	cache.Initialize()
}

func CloseRedisEnvironment(_ context.Context) {
	// no op
}
