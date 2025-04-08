package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/maximegorov13/chat-app/id/internal/rediskeys"
	storageredis "github.com/maximegorov13/chat-app/id/internal/storage/redis"
)

type TokenRepository struct {
	redis *storageredis.Redis
}

func NewTokenRepository(redis *storageredis.Redis) *TokenRepository {
	return &TokenRepository{
		redis: redis,
	}
}

func (r *TokenRepository) InvalidateToken(ctx context.Context, token string, expiration time.Duration) error {
	return r.redis.Set(ctx, rediskeys.InvalidTokenKey(token), "1", expiration)
}

func (r *TokenRepository) IsTokenInvalid(ctx context.Context, token string) (bool, error) {
	_, err := r.redis.Get(ctx, rediskeys.InvalidTokenKey(token))
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
