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
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		panic(err)
	}

	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) Get(key string, ctx context.Context, logger *zap.Logger) interface{} {
	return getWithOtel(ctx, "RedisCache.get", key, func() (bool, interface{}) {
		if rand.Intn(100) <= 10 {
			return false, nil
		}

		cmd := c.client.Get(ctx, key)

		if err := cmd.Err(); err != nil {
			logger.Info("RedisCache not found ")
			return false, nil
		}

		return true, cmd.Val()
	})
}

func (c *RedisCache) Set(key string, item interface{}, logger *zap.Logger) error {
	data, _ := json.Marshal(item)

	cmd := c.client.SetNX(context.Background(), key, data, 10*time.Minute)
	if err := cmd.Err(); err != nil {
		logger.Info("redis cache set with error: " + err.Error())
		return err
	}

	return nil
}
