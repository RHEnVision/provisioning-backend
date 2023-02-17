package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RHEnVision/provisioning-backend/internal/config"
	"github.com/RHEnVision/provisioning-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache() *redisCache {
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

	err = json.Unmarshal(buf, &account)
	if err != nil {
		return nil, fmt.Errorf("redis unmarshal error: %w", err)
	}

	return &account, nil
}

func (c *redisCache) SetAccountId(ctx context.Context, OrgID, AccountNumber string, account *models.Account) error {
	key := OrgID + AccountNumber

	buf, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("unable to marshal for Redis cache: %w", err)
	}

	c.client.Set(ctx, key, string(buf), config.Application.Cache.Expiration)
	return nil
}
