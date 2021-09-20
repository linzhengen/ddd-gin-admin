package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type RoleMenuRepository interface {
	Query(ctx context.Context, params schema.RoleMenuQueryParam) (entity.RoleMenus, *schema.PaginationResult, error)
	Get(ctx context.Context, id string) (*entity.RoleMenu, error)
	Create(ctx context.Context, item entity.RoleMenu) error
	Update(ctx context.Context, id string, item entity.RoleMenu) error
	Delete(ctx context.Context, id string) error
	DeleteByRoleID(ctx context.Context, roleID string) error
}
