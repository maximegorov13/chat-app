package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

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
	return r.redis.Set(ctx, fmt.Sprintf("invalid_token:%s", token), "1", expiration)
}

func (r *TokenRepository) IsTokenInvalid(ctx context.Context, token string) (bool, error) {
	_, err := r.redis.Get(ctx, fmt.Sprintf("invalid_token:%s", token))
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
