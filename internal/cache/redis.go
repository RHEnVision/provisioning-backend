// Package cache provides application cache based on Redis. This feature can be turned off
// via configuration and in that case function Find return ErrNotFound and functions Set
// do nothing.
package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/RHEnVision/provisioning-backend/internal/clients"
	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/metrics"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrNotFound = errors.New("not found in cache")
	ErrNilValue = errors.New("value is nil")

	// when redis is enabled via configuration
	redisEnabled bool

	// the client
	client *redis.Client

	// application id "constant" memory-only cache
	appTypeId      *string
	appTypeIdMutex sync.Mutex
)

// Forever is used for items that should be cached "forever". Expiration of 30 days
// is used to allow cleanup of unused items.
const Forever time.Duration = 24 * time.Hour * 30

type Cacheable interface {
	CacheKeyName() string
}

// Initialize creates new Redis client if allowed by application config, or does nothing.
func Initialize() {
	if config.Application.Cache.Type == "redis" {
		redisEnabled = true
		log.Logger.Info().Bool("cache", true).Msg("Initializing redis application cache")

		// register all Cacheable types
		gob.Register(&models.Account{})
		gob.Register(&clients.AccountDetailsAWS{})
		gob.Register(&clients.AccessList{})

		client = redis.NewClient(&redis.Options{
			Addr:     config.RedisHostAndPort(),
			Username: config.Application.Cache.Redis.User,
			Password: config.Application.Cache.Redis.Password,
			DB:       config.Application.Cache.Redis.DB,
		})
	} else {
		log.Logger.Info().Bool("cache", true).Msg("No application cache in use")
	}
}

// FindAppTypeId returns "application id" special identifier, or returns ErrNotFound.
func FindAppTypeId(_ context.Context) (string, error) {
	appTypeIdMutex.Lock()
	defer appTypeIdMutex.Unlock()

	if appTypeId == nil {
		return "", ErrNotFound
	}
	return *appTypeId, nil
}

// SetAppTypeId sets "application id" special identifier.
func SetAppTypeId(_ context.Context, value string) error {
	appTypeIdMutex.Lock()
	defer appTypeIdMutex.Unlock()

	appTypeId = &value
	return nil
}

// Find returns an item from cache. ErrNotFound is returned on cache miss or when
// the item cannot be deserialized
func Find(ctx context.Context, key string, value Cacheable) error {
	if !redisEnabled {
		return ErrNotFound
	}

	if value == nil {
		return ErrNilValue
	}

	prefix := value.CacheKeyName()

	cmd := client.Get(ctx, prefix+key)
	if errors.Is(cmd.Err(), redis.Nil) {
		metrics.IncCacheHit(prefix, "miss")
		return ErrNotFound
	} else if cmd.Err() != nil {
		metrics.IncCacheHit(prefix, "err")
		return fmt.Errorf("redis get error: %w", cmd.Err())
	}

	buf, err := cmd.Bytes()
	if err != nil {
		metrics.IncCacheHit(prefix, "err")
		return fmt.Errorf("redis bytes conversion error: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(buf))

	err = dec.Decode(value)
	if err != nil {
		// decode error can be thrown if previous cache entry was JSON-encoded, return not found to overwrite it
		zerolog.Ctx(ctx).Warn().Err(err).Bool("cache", true).Msgf("Redis cache decode error: %s", err.Error())
		metrics.IncCacheHit(prefix, "err")
		return ErrNotFound
	}

	metrics.IncCacheHit(prefix, "hit")
	zerolog.Ctx(ctx).Trace().Bool("cache", true).Msgf("Cache hit for key '%s%s' type %T", prefix, key, value)
	return nil
}

// SetExpires calls Set with specific expiration.
func SetExpires(ctx context.Context, key string, value Cacheable, expiration time.Duration) error {
	if !redisEnabled {
		return nil
	}

	if value == nil {
		return ErrNilValue
	}

	prefix := value.CacheKeyName()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(value)
	if err != nil {
		metrics.IncCacheHit(prefix, "err")
		return fmt.Errorf("unable to encode for Redis cache: %w", err)
	}

	cmd := client.Set(ctx, prefix+key, buf.String(), expiration)
	if cmd.Err() != nil {
		metrics.IncCacheHit(prefix, "err")
		return fmt.Errorf("redis set error: %w", cmd.Err())
	}

	return nil
}

// SetForever calls Set with Forever expiration duration.
// nolint: wrapcheck
func SetForever(ctx context.Context, key string, value Cacheable) error {
	return SetExpires(ctx, key, value, Forever)
}

// Set creates or updates existing cache entry. It uses the default expiration duration
// specified in the application configuration.
// nolint: wrapcheck
func Set(ctx context.Context, key string, value Cacheable) error {
	return SetExpires(ctx, key, value, config.Application.Cache.Expiration)
}
