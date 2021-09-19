package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/response"

	"github.com/linzhengen/ddd-gin-admin/app/domain/entity"

	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/request"
)

type MenuActionResourceRepository interface {
	Query(ctx context.Context, req request.MenuActionResourceQuery) (entity.MenuActionResources, *response.Pagination, error)
	Get(ctx context.Context, id string) (*entity.MenuActionResource, error)
	Create(ctx context.Context, item entity.MenuActionResource) error
	Update(ctx context.Context, id string, item entity.MenuActionResource) error
	Delete(ctx context.Context, id string) error
	DeleteByActionID(ctx context.Context, actionID string) error
	DeleteByMenuID(ctx context.Context, menuID string) error
}
