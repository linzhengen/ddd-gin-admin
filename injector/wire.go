//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"
	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/casbin"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	consoleHandler "github.com/linzhengen/ddd-gin-admin/app/interfaces/console/handler"
	"github.com/linzhengen/ddd-gin-admin/injector/api"
	"github.com/linzhengen/ddd-gin-admin/injector/console"

	// "github.com/linzhengen/ddd-gin-admin/infrastructure/api/mock"
	"github.com/google/wire"
)

func BuildApiInjector() (*ApiInjector, func(), error) {
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

		// factory
		factory.NewMenu,
		factory.NewMenuAction,
		factory.NewMenuActionResource,
		factory.NewRole,
		factory.NewRoleMenu,
		factory.NewUser,
		factory.NewUserRole,

		// service
		service.NewLogin,
		service.NewMenu,
		service.NewRole,
		service.NewUser,

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
		NewApiInjector,
	)
	return nil, nil, nil
}

func BuildConsoleInjector() (*consoleHandler.ConsoleHandler, error) {
	wire.Build(
		console.InitGormDB,
		persistence.NewUser,
		application.NewHelloConsole,
		consoleHandler.NewHelloHandler,
		consoleHandler.NewConsoleHandler,
	)
	return nil, nil
}
