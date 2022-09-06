package cache

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisCache struct {
	logger *zap.Logger
	client *redis.Client
}

func NewRedisCache(logger *zap.Logger) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		DB:   0,
	})

	return &RedisCache{
		logger: logger,
		client: client,
	}
}

func (c *RedisCache) Get(key string, ctx context.Context) interface{} {
	return getWithOpel(ctx, "RedisCache.get", key, func() (bool, interface{}) {
		if rand.Intn(100) < 10 {
			return false, nil
		}

		// random sleep time for mock
		time.Sleep(time.Duration(rand.Intn(30)) * time.Millisecond)

		cmd := c.client.Get(ctx, key)
		if err := cmd.Err(); err != nil {
			c.logger.Info("RedisCache not found " + key)
			return false, nil
		}

		return true, cmd.Val()
	})
}

func (c *RedisCache) Set(key string, item interface{}) error {
	data, _ := json.Marshal(item)

	cmd := c.client.SetNX(context.Background(), key, data, 10*time.Minute)
	if err := cmd.Err(); err != nil {
		c.logger.Error("RedisCache.set with error: " + err.Error())
		return err
	}

	return nil
}