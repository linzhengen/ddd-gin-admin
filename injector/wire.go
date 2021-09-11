//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/casbin"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/persistence"
	"github.com/linzhengen/ddd-gin-admin/injector/api"
	"github.com/linzhengen/ddd-gin-admin/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/interfaces/api/router"

	// "github.com/linzhengen/ddd-gin-admin/infrastructure/api/mock"
	"github.com/google/wire"
)

func BuildApiInjector() (*api.Injector, func(), error) {
	wire.Build(
		// init,
		api.InitGormDB,
		api.InitAuth,
		api.InitCasbin,
		api.InitGinEngine,

		// persistence
		persistence.NewTrans,
		persistence.NewUser,
		persistence.NewRole,
		persistence.NewUserRole,
		persistence.NewMenu,
		persistence.NewRoleMenu,
		persistence.NewMenuAction,
		persistence.NewMenuActionResource,

		// application
		application.NewLogin,
		application.NewMenu,
		application.NewRole,
		application.NewUser,

		// handler
		handler.NewMenu,
		handler.NewRole,
		handler.NewLogin,
		handler.NewUser,
		handler.NewHealthCheck,

		// router
		router.NewRouter,

		// lib
		casbin.NewCasbinAdapter,

		// injector
		api.NewInjector,
	)
	return nil, nil, nil
}
