// Package mem contains an in-memory cache implementation that is thread-safe with unlimited capacity.
package mem

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/cache"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	ttl "github.com/akyoto/cache"
)

type memoryCache struct {
	appTypeIdCache *string
	appTypeIdMutex sync.Mutex

	instanceTypeCache *ttl.Cache
	instanceTypeMutex sync.Mutex
}

var memCache memoryCache

const CleanupInterval = 5 * time.Minute

func init() {
	memCache = memoryCache{
		instanceTypeCache: ttl.New(CleanupInterval),
	}

	cache.GetGlobalCache = getMemoryGlobalCache
}

func getMemoryGlobalCache() cache.Cache {
	return &memCache
}

func logCacheMiss(ctx context.Context, modelName string) {
	logger := ctxval.Logger(ctx)
	logger.Trace().Str("cache", modelName).Msg("Cache miss")
}

func (c *memoryCache) AppTypeId(ctx context.Context) (string, bool) {
	c.appTypeIdMutex.Lock()
	defer c.appTypeIdMutex.Unlock()

	if c.appTypeIdCache == nil {
		logCacheMiss(ctx, "appTypeId")
		return "", false
	}
	return *c.appTypeIdCache, true
}

func (c *memoryCache) SetAppTypeId(_ context.Context, appTypeId string) {
	c.appTypeIdMutex.Lock()
	defer c.appTypeIdMutex.Unlock()

	c.appTypeIdCache = &appTypeId
}

func cacheKey(tokens ...string) string {
	return strings.Join(tokens, "-")
}

func (c *memoryCache) InstanceTypes(ctx context.Context, sourceId, awsRegion string) ([]string, bool) {
	c.instanceTypeMutex.Lock()
	defer c.instanceTypeMutex.Unlock()

	value, ok := c.instanceTypeCache.Get(cacheKey("it", sourceId, awsRegion))
	if !ok {
		logCacheMiss(ctx, "instanceType")
	}
	return value.([]string), ok
}

func (c *memoryCache) SetInstanceTypes(_ context.Context, sourceId, awsRegion string, types []string, ttl time.Duration) {
	c.instanceTypeMutex.Lock()
	defer c.instanceTypeMutex.Unlock()

	c.instanceTypeCache.Set(cacheKey("it", sourceId, awsRegion), types, ttl)
}
