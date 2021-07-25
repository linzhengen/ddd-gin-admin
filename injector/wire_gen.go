// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/domain/adapter"
	"github.com/linzhengen/ddd-gin-admin/domain/repository"
	"github.com/linzhengen/ddd-gin-admin/interfaces/handler"
	"github.com/linzhengen/ddd-gin-admin/interfaces/router"
)

import (
	_ "github.com/linzhengen/ddd-gin-admin/interfaces/swagger"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, func(), error) {
	author, cleanup, err := InitAuth()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := InitGormDB()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	role := &repository.Role{
		DB: db,
	}
	roleMenu := &repository.RoleMenu{
		DB: db,
	}
	menuActionResource := &repository.MenuActionResource{
		DB: db,
	}
	user := &repository.User{
		DB: db,
	}
	userRole := &repository.UserRole{
		DB: db,
	}
	casbinAdapter := &adapter.CasbinAdapter{
		RoleModel:         role,
		RoleMenuModel:     roleMenu,
		MenuResourceModel: menuActionResource,
		UserModel:         user,
		UserRoleModel:     userRole,
	}
	syncedEnforcer, cleanup3, err := InitCasbin(casbinAdapter)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	demo := &repository.Demo{
		DB: db,
	}
	applicationDemo := &application.Demo{
		DemoModel: demo,
	}
	handlerDemo := &handler.Demo{
		DemoSrv: applicationDemo,
	}
	menu := &repository.Menu{
		DB: db,
	}
	menuAction := &repository.MenuAction{
		DB: db,
	}
	login := &application.Login{
		Auth:            author,
		UserModel:       user,
		UserRoleModel:   userRole,
		RoleModel:       role,
		RoleMenuModel:   roleMenu,
		MenuModel:       menu,
		MenuActionModel: menuAction,
	}
	handlerLogin := &handler.Login{
		LoginSrv: login,
	}
	trans := &repository.Trans{
		DB: db,
	}
	applicationMenu := &application.Menu{
		TransModel:              trans,
		MenuModel:               menu,
		MenuActionModel:         menuAction,
		MenuActionResourceModel: menuActionResource,
	}
	handlerMenu := &handler.Menu{
		MenuSrv: applicationMenu,
	}
	applicationRole := &application.Role{
		Enforcer:      syncedEnforcer,
		TransModel:    trans,
		RoleModel:     role,
		RoleMenuModel: roleMenu,
		UserModel:     user,
	}
	handlerRole := &handler.Role{
		RoleSrv: applicationRole,
	}
	applicationUser := &application.User{
		Enforcer:      syncedEnforcer,
		TransModel:    trans,
		UserModel:     user,
		UserRoleModel: userRole,
		RoleModel:     role,
	}
	handlerUser := &handler.User{
		UserSrv: applicationUser,
	}
	routerRouter := &router.Router{
		Auth:           author,
		CasbinEnforcer: syncedEnforcer,
		DemoAPI:        handlerDemo,
		LoginAPI:       handlerLogin,
		MenuAPI:        handlerMenu,
		RoleAPI:        handlerRole,
		UserAPI:        handlerUser,
	}
	engine := InitGinEngine(routerRouter)
	injector := &Injector{
		Engine:         engine,
		Auth:           author,
		CasbinEnforcer: syncedEnforcer,
		MenuBll:        applicationMenu,
	}
	return injector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
