package application

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
)

type DbSeedConsole interface {
	Seed(ctx context.Context) error
}

func NewDbSeedConsole(menuSvc service.Menu) DbSeedConsole {
	return &dbSeedConsole{
		menuSvc: menuSvc,
	}
}

type dbSeedConsole struct {
	menuSvc service.Menu
}

func (d dbSeedConsole) Seed(ctx context.Context) error {
	return d.menuSvc.InitData(ctx, "./configs/menu.yaml")
}
