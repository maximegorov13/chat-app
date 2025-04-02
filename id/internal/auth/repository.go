package auth

import (
	"context"
	"time"
)

type TokenRepository interface {
	InvalidateToken(ctx context.Context, token string, expiration time.Duration) error
	IsTokenInvalid(ctx context.Context, token string) (bool, error)
}
