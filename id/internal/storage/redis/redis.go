package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/maximegorov13/chat-app/id/configs"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, conf *configs.Config) (*Redis, error) {
	opts, err := redis.ParseURL(conf.Redis.Url)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)
	if err = rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{
		client: rdb,
	}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
