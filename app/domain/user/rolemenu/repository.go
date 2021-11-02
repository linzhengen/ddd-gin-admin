package rolemenu

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Repository interface {
	Query(ctx context.Context, params QueryParam) (RoleMenus, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*RoleMenu, error)
	Create(ctx context.Context, item *RoleMenu) error
	Update(ctx context.Context, id string, item *RoleMenu) error
	Delete(ctx context.Context, id string) error
	DeleteByRoleID(ctx context.Context, roleID string) error
}
