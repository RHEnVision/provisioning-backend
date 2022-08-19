// Package noop implements no operation cache, a cache that does not store any data
// and always misses. It is used in tests.
package noop

import (
	"context"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
)

type noopCache struct{}

var noCache noopCache

func init() {
	noCache = noopCache{}
	cache.GetGlobalCache = getNoopCache
}

func getNoopCache() cache.Cache {
	return &noCache
}

func (c *noopCache) AppTypeId(ctx context.Context) (string, bool) {
	return "", false
}

func (c *noopCache) SetAppTypeId(_ context.Context, appTypeId string) {
}

func (c *noopCache) InstanceTypes(ctx context.Context, sourceId, awsRegion string) ([]string, bool) {
	return nil, false
}

func (c *noopCache) SetInstanceTypes(_ context.Context, sourceId, awsRegion string, types []string, ttl time.Duration) {
}
