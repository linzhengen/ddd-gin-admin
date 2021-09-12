package repository

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/schema"
)

type MenuActionRepository interface {
	Query(ctx context.Context, params schema.MenuActionQueryParam, opts ...schema.MenuActionQueryOptions) (*schema.MenuActionQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.MenuActionQueryOptions) (*schema.MenuAction, error)
	Create(ctx context.Context, item schema.MenuAction) error
	Update(ctx context.Context, id string, item schema.MenuAction) error
	Delete(ctx context.Context, id string) error
	DeleteByMenuID(ctx context.Context, menuID string) error
}
