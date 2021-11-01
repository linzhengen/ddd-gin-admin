package auth

import (
	"context"
)

type Repository interface {
	FindRootUser(ctx context.Context, userName string) *RootUser
	GenerateToken(ctx context.Context, userID string) (*Auth, error)
	DestroyToken(ctx context.Context, accessToken string) error
	ParseUserID(ctx context.Context, accessToken string) (string, error)
	Release() error
}
