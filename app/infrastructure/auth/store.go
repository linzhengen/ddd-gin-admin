package auth

import (
	"context"
	"time"
)

type Store interface {
	Set(ctx context.Context, tokenString string, expiration time.Duration) error
	Check(ctx context.Context, tokenString string) (bool, error)
	Close() error
}
