package cache

import (
	"context"
	"time"
)

// Cache provides an application global thread-safe cache.
type Cache interface {
	AppTypeIdCache
	InstanceTypesCache
}

var GetGlobalCache func() Cache

type AppTypeIdCache interface {
	AppTypeId(ctx context.Context) (string, bool)
	SetAppTypeId(ctx context.Context, appTypeId string)
}

type InstanceTypesCache interface {
	InstanceTypes(ctx context.Context, sourceId, awsRegion string) ([]string, bool)
	SetInstanceTypes(_ context.Context, sourceId, awsRegion string, types []string, ttl time.Duration)
}
