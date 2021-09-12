package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
)

type UserRepository interface {
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error)
	Create(ctx context.Context, item schema.User) error
	Update(ctx context.Context, id string, item schema.User) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
	UpdatePassword(ctx context.Context, id, password string) error
}
