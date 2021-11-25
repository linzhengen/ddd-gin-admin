// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/app/application"
	menu2 "github.com/linzhengen/ddd-gin-admin/app/domain/menu"
	user2 "github.com/linzhengen/ddd-gin-admin/app/domain/user"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuaction"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/menu/menuactionresource"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/trans"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/role"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/rolemenu"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/user/userrole"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/injector/api"
)

import (
	_ "github.com/linzhengen/ddd-gin-admin/app/interfaces/api/swagger"
)

// Injectors from wire.go:

func BuildApiInjector() (*ApiInjector, func(), error) {
	repository, cleanup, err := api.InitAuth()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := api.InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userRepository := user.NewRepository(db)
	roleRepository := role.NewRepository(db)
	userroleRepository := userrole.NewRepository(db)
	service := user2.NewService(repository, userRepository, roleRepository, userroleRepository)
	menuRepository := menu.NewRepository(db)
	menuactionRepository := menuaction.NewRepository(db)
	rolemenuRepository := rolemenu.NewRepository(db)
	login := application.NewLogin(repository, userRepository, roleRepository, userroleRepository, service, menuRepository, menuactionRepository, rolemenuRepository)
	handlerLogin := handler.NewLogin(login)
	transRepository := trans.NewRepository(db)
	menuactionresourceRepository := menuactionresource.NewRepository(db)
	menuService := menu2.NewService(transRepository, menuRepository, menuactionRepository, menuactionresourceRepository)
	applicationMenu := application.NewMenu(menuService)
	handlerMenu := handler.NewMenu(applicationMenu)
	applicationRole := application.NewRole(transRepository, roleRepository, rolemenuRepository, userRepository)
	handlerRole := handler.NewRole(applicationRole)
	applicationUser := application.NewUser(repository, transRepository, userRepository, userroleRepository, roleRepository)
	handlerUser := handler.NewUser(applicationUser)
	healthCheck := handler.NewHealthCheck()
	routerRouter := router.NewRouter(repository, handlerLogin, handlerMenu, handlerRole, handlerUser, healthCheck)
	engine := api.InitGinEngine(routerRouter)
	apiInjector := NewApiInjector(engine, repository)
	return apiInjector, func() {
		cleanup2()
		cleanup()
	}, nil
}
