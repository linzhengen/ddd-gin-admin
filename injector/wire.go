//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/trans"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/injector/api"

	"github.com/google/wire"
)

func BuildApiInjector() (*ApiInjector, func(), error) {
	wire.Build(
		// init,
		api.InitGormDB,
		api.InitAuth,
		api.InitGinEngine,

		// infrastructure
		menu.NewRepository,
		menuaction.NewRepository,
		menuactionresource.NewRepository,
		user.NewRepository,
		userrole.NewRepository,
		rolemenu.NewRepository,
		role.NewRepository,
		trans.NewRepository,
		//auth.NewRepository,

		// application
		application.NewMenu,
		application.NewRole,
		application.NewUser,
		application.NewLogin,

		// router
		router.NewRouter,

		// injector
		NewApiInjector,
	)
	return nil, nil, nil
}
