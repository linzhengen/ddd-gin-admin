package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"
)

type Menu interface {
	Query(ctx context.Context, params menu.QueryParam) (menu.Menus, *pagination.Pagination, error)
	Get(ctx context.Context, id string) (*menu.Menu, error)
	QueryActions(ctx context.Context, id string) (menuaction.MenuActions, error)
	Create(ctx context.Context, item *menu.Menu) (string, error)
	Update(ctx context.Context, id string, item *menu.Menu) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status int) error
}

func NewMenu(
	menuSvc menu.Service,
) Menu {
	return &menuApp{
		menuSvc: menuSvc,
	}
}

type menuApp struct {
	menuSvc menu.Service
}

func (m menuApp) Query(ctx context.Context, params menu.QueryParam) (menu.Menus, *pagination.Pagination, error) {
	return m.menuSvc.Query(ctx, params)
}

func (m menuApp) Get(ctx context.Context, id string) (*menu.Menu, error) {
	return m.menuSvc.Get(ctx, id)
}

func (m menuApp) QueryActions(ctx context.Context, id string) (menuaction.MenuActions, error) {
	return m.menuSvc.QueryActions(ctx, id)
}

func (m menuApp) Create(ctx context.Context, item *menu.Menu) (string, error) {
	return m.menuSvc.Create(ctx, item)
}

func (m menuApp) Update(ctx context.Context, id string, item *menu.Menu) error {
	return m.menuSvc.Update(ctx, id, item)
}

func (m menuApp) Delete(ctx context.Context, id string) error {
	return m.menuSvc.Delete(ctx, id)
}

func (m menuApp) UpdateStatus(ctx context.Context, id string, status int) error {
	return m.menuSvc.UpdateStatus(ctx, id, status)
}
