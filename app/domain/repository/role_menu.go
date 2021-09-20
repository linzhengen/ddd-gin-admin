package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type RoleMenuRepository interface {
	Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error)
	Create(ctx context.Context, item schema.RoleMenu) error
	Update(ctx context.Context, id string, item schema.RoleMenu) error
	Delete(ctx context.Context, id string) error
	DeleteByRoleID(ctx context.Context, roleID string) error
}
