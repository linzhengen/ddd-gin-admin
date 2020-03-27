package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/domain/schema"
)

// UserRepository is interface.
type UserRepository interface {
	Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error)
	QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error)
	Get(ctx context.Context, recordID string, opts ...schema.UserQueryOptions) (*schema.User, error)
	Create(ctx context.Context, item schema.User) (*schema.User, error)
	Update(ctx context.Context, recordID string, item schema.User) (*schema.User, error)
	Delete(ctx context.Context, recordID string) error
	UpdateStatus(ctx context.Context, recordID string, status int) error
}
