package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/ctxval"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() AccountIdCache {
	// register gob types
	gob.Register(&models.Account{})

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisHostAndPort(),
		Username: config.Application.Cache.Redis.User,
		Password: config.Application.Cache.Redis.Password,
		DB:       config.Application.Cache.Redis.DB,
	})
	return &redisCache{
		client: client,
	}
}

func (c *redisCache) FindAccountId(ctx context.Context, OrgID, AccountNumber string) (*models.Account, error) {
	key := OrgID + AccountNumber
	account := models.Account{}

	cmd := c.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redis.Nil) {
		return nil, NotFound
	} else if cmd.Err() != nil {
		return nil, fmt.Errorf("redis error: %w", cmd.Err())
	}

	buf, err := cmd.Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis bytes conversion error: %w", err)
	}

	dec := gob.NewDecoder(bytes.NewReader(buf))

	err = dec.Decode(&account)
	if err != nil {
		// decode error can be thrown if previous cache entry was JSON-encoded, return not found to overwrite it
		ctxval.Logger(ctx).Warn().Err(err).Bool("cache", true).Msgf("redis cache decode error: %s", err.Error())
		return nil, NotFound
	}

	return &account, nil
}

func (c *redisCache) SetAccountId(ctx context.Context, OrgID, AccountNumber string, account *models.Account) error {
	key := OrgID + AccountNumber

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(account)
	if err != nil {
		return fmt.Errorf("unable to encode for Redis cache: %w", err)
	}

	c.client.Set(ctx, key, buf.String(), config.Application.Cache.Expiration)
	return nil
}
