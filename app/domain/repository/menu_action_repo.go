package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"
)

type MenuActionRepository interface {
	Query(ctx context.Context, req request.MenuActionQuery) (entity.MenuActions, *response.Pagination, error)
	Get(ctx context.Context, id string) (*entity.MenuAction, error)
	Create(ctx context.Context, item entity.MenuAction) error
	Update(ctx context.Context, id string, item entity.MenuAction) error
	Delete(ctx context.Context, id string) error
	DeleteByMenuID(ctx context.Context, menuID string) error
}
