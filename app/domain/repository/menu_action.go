package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type MenuActionRepository interface {
	Query(ctx context.Context, params schema.MenuActionQueryParam) (entity.MenuActions, *schema.PaginationResult, error)
	Get(ctx context.Context, id string) (*entity.MenuAction, error)
	Create(ctx context.Context, item entity.MenuAction) error
	Update(ctx context.Context, id string, item entity.MenuAction) error
	Delete(ctx context.Context, id string) error
	DeleteByMenuID(ctx context.Context, menuID string) error
}
