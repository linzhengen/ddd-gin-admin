package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/schema"
)

type Menu interface {
	InitData(ctx context.Context, dataFile string) error
	Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error)
	Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error)
	QueryActions(ctx context.Context, id string) (schema.MenuActions, error)
	Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error)
	Update(ctx context.Context, id string, item schema.Menu) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewMenu(
	menuSvc service.Menu,
) Menu {
	return &menu{
		menuSvc: menuSvc,
	}
}

type menu struct {
	menuSvc service.Menu
}

func (m menu) InitData(ctx context.Context, dataFile string) error {
	return m.menuSvc.InitData(ctx, dataFile)
}

func (m menu) Query(ctx context.Context, params schema.MenuQueryParam, opts ...schema.MenuQueryOptions) (*schema.MenuQueryResult, error) {
	return m.menuSvc.Query(ctx, params, opts...)
}

func (m menu) Get(ctx context.Context, id string, opts ...schema.MenuQueryOptions) (*schema.Menu, error) {
	return m.menuSvc.Get(ctx, id, opts...)
}

func (m menu) QueryActions(ctx context.Context, id string) (schema.MenuActions, error) {
	return m.menuSvc.QueryActions(ctx, id)
}

func (m menu) Create(ctx context.Context, item schema.Menu) (*schema.IDResult, error) {
	panic("implement me")
}

func (m menu) Update(ctx context.Context, id string, item schema.Menu) error {
	panic("implement me")
}

func (m menu) Delete(ctx context.Context, id string) error {
	panic("implement me")
}

func (m menu) UpdateStatus(ctx context.Context, id string, status int) error {
	panic("implement me")
}
