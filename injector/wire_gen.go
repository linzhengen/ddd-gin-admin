// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"context"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure"

	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/role"

	"github.com/linzhengen/ddd-gin-admin/app/application"
	"github.com/linzhengen/ddd-gin-admin/app/domain/factory"
	"github.com/linzhengen/ddd-gin-admin/app/domain/service"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/casbin"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/persistence"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/console/command"
	"github.com/linzhengen/ddd-gin-admin/injector/api"
	"github.com/linzhengen/ddd-gin-admin/injector/console"

	_ "github.com/linzhengen/ddd-gin-admin/app/interfaces/api/swagger"
)

// Injectors from wire.go:

func BuildApiInjector() (*ApiInjector, func(), error) {
	author, cleanup, err := api.InitAuth()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := api.InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	roleRepository := role.NewRole(db)
	roleMenuRepository := persistence.NewRoleMenu(db)
	menuActionResourceRepository := persistence.NewMenuActionResource(db)
	userRepository := user.NewUser(db)
	userRoleRepository := persistence.NewUserRole(db)
	role := factory.NewRole()
	roleMenu := factory.NewRoleMenu()
	menuActionResource := factory.NewMenuActionResource()
	user := factory.NewUser()
	userRole := factory.NewUserRole()
	casbinAdapter := casbin.NewCasbinAdapter(roleRepository, roleMenuRepository, menuActionResourceRepository, userRepository, userRoleRepository, role, roleMenu, menuActionResource, user, userRole)
	syncedEnforcer, cleanup3, err := api.InitCasbin(casbinAdapter)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	menuRepository := persistence.NewMenu(db)
	menuActionRepository := persistence.NewMenuAction(db)
	menuAction := factory.NewMenuAction()
	login := service.NewLogin(author, userRepository, userRoleRepository, roleRepository, roleMenuRepository, menuRepository, menuActionRepository, user, userRole, role, roleMenu, menuAction)
	applicationLogin := application.NewLogin(login)
	handlerLogin := handler.NewLogin(applicationLogin)
	transRepository := persistence.NewTrans(db)
	menu := factory.NewMenu()
	serviceMenu := service.NewMenu(transRepository, menuRepository, menuActionRepository, menuActionResourceRepository, menu, menuAction, menuActionResource)
	applicationMenu := application.NewMenu(serviceMenu)
	handlerMenu := handler.NewMenu(applicationMenu)
	serviceRole := service.NewRole(casbinAdapter, syncedEnforcer, transRepository, roleRepository, roleMenuRepository, userRepository, role, roleMenu, user)
	applicationRole := application.NewRole(serviceRole)
	handlerRole := handler.NewRole(applicationRole)
	serviceUser := service.NewUser(casbinAdapter, syncedEnforcer, transRepository, userRepository, userRoleRepository, roleRepository, user, userRole, role)
	applicationUser := application.NewUser(serviceUser)
	handlerUser := handler.NewUser(applicationUser)
	healthCheck := handler.NewHealthCheck()
	routerRouter := router.NewRouter(author, syncedEnforcer, handlerLogin, handlerMenu, handlerRole, handlerUser, healthCheck)
	engine := api.InitGinEngine(routerRouter)
	apiInjector := NewApiInjector(engine, author, syncedEnforcer, serviceMenu)
	return apiInjector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

func BuildConsoleInjector(ctx context.Context) (command.Commands, func(), error) {
	db, cleanup, err := console.InitGormDB()
	if err != nil {
		return nil, nil, err
	}
	dbMigrationRepository := infrastructure.NewDbMigration(db)
	dbMigrationConsole := application.NewDbMigrationConsole(dbMigrationRepository)
	migrateCommand := command.NewMigrateCommand(dbMigrationConsole)
	transRepository := persistence.NewTrans(db)
	menuRepository := persistence.NewMenu(db)
	menuActionRepository := persistence.NewMenuAction(db)
	menuActionResourceRepository := persistence.NewMenuActionResource(db)
	menu := factory.NewMenu()
	menuAction := factory.NewMenuAction()
	menuActionResource := factory.NewMenuActionResource()
	serviceMenu := service.NewMenu(transRepository, menuRepository, menuActionRepository, menuActionResourceRepository, menu, menuAction, menuActionResource)
	dbSeedConsole := application.NewDbSeedConsole(serviceMenu)
	seedCommand := command.NewSeedCommand(dbSeedConsole)
	commands := command.NewCliCommands(ctx, migrateCommand, seedCommand)
	return commands, func() {
		cleanup()
	}, nil
}
