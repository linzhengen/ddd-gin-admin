//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/menu"
	"github.com/linzhengen/ddd-gin-admin/app/domain/user"
	menuInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu"
	menuActionInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuaction"
	menuActionResourceInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuactionresource"

	rbacInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/rbac"
	transInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/trans"
	userInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"
	roleInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/role"
	roleMenuInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/rolemenu"
	userRoleInfra "github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/injector/api"

	"github.com/google/wire"
)

func BuildApiInjector() (*ApiInjector, func(), error) {
	wire.Build(
		// init,
		InitGormDB,
		api.InitAuth,
		api.InitGinEngine,
		api.InitCasbin,

		// domain
		user.NewService,
		menu.NewService,

		// infrastructure
		menuInfra.NewRepository,
		menuActionInfra.NewRepository,
		menuActionResourceInfra.NewRepository,
		userInfra.NewRepository,
		userRoleInfra.NewRepository,
		roleMenuInfra.NewRepository,
		roleInfra.NewRepository,
		transInfra.NewRepository,
		rbacInfra.NewRepository,

		// application
		application.NewMenu,
		application.NewRole,
		application.NewUser,
		application.NewLogin,
		application.NewRbacAdapter,
		application.NewSeed,

		// handler
		handler.NewHealthCheck,
		handler.NewUser,
		handler.NewRole,
		handler.NewMenu,
		handler.NewLogin,

		// router
		router.NewRouter,

		// injector
		NewApiInjector,
	)
	return nil, nil, nil
}
