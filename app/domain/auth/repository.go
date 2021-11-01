package auth

import (
	"context"
)

type Repository interface {
	GenerateToken(ctx context.Context, userID string) (*Auth, error)
	DestroyToken(ctx context.Context, accessToken string) error
	ParseUserID(ctx context.Context, accessToken string) (string, error)
	Release() error
}
