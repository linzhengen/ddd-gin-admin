package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/domain/schema"
)

type RoleRepository interface {
	Query(ctx context.Context, params schema.RoleQueryParam, opts ...schema.RoleQueryOptions) (*schema.RoleQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.RoleQueryOptions) (*schema.Role, error)
	Create(ctx context.Context, item schema.Role) error
	Update(ctx context.Context, id string, item schema.Role) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}
