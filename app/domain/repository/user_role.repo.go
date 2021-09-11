package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
)

type UserRoleRepository interface {
	Query(ctx context.Context, params schema.UserRoleQueryParam, opts ...schema.UserRoleQueryOptions) (*schema.UserRoleQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.UserRoleQueryOptions) (*schema.UserRole, error)
	Create(ctx context.Context, item schema.UserRole) error
	Update(ctx context.Context, id string, item schema.UserRole) error
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
