package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streamverse/common-go/logger"
	"go.uber.org/zap"
)

type RedisClient struct {
	client *redis.Client
	log    *logger.Logger
}

func NewRedisClient(addr, password string, db int, log *logger.Logger) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Error("Failed to connect to Redis", zap.Error(err))
	} else {
		log.Info("Connected to Redis", zap.String("addr", addr))
	}

	return &RedisClient{
		client: client,
		log:    log,
	}
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, jsonVal, expiration).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
