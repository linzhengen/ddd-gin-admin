package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type MenuRepository interface {
	Query(ctx context.Context, params schema.MenuQueryParam) (entity.Menus, *schema.PaginationResult, error)
	Get(ctx context.Context, id string) (*entity.Menu, error)
	Create(ctx context.Context, item entity.Menu) error
	Update(ctx context.Context, id string, item entity.Menu) error
	UpdateParentPath(ctx context.Context, id, parentPath string) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}
