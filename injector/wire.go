//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"
	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/casbin"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/console/command"
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

		// command
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

func BuildConsoleInjector(ctx context.Context) (command.Commands, func(), error) {
	wire.Build(
		// init
		console.InitGormDB,

		// factory
		factory.NewMenu,
		factory.NewMenuAction,
		factory.NewMenuActionResource,

		// persistence
		persistence.NewDbMigration,
		persistence.NewTrans,
		persistence.NewMenu,
		persistence.NewMenuAction,
		persistence.NewMenuActionResource,

		// service
		service.NewMenu,

		// application
		application.NewDbMigrationConsole,
		application.NewDbSeedConsole,

		// command
		command.NewMigrateCommand,
		command.NewSeedCommand,
		command.NewCliCommands,
	)
	return nil, nil, nil
}
