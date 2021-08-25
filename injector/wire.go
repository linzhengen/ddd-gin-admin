//go:build wireinject

// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/injector/api"
	"github.com/linzhengen/ddd-gin-admin/interfaces/api/handler"
	"github.com/linzhengen/ddd-gin-admin/interfaces/api/router"

	// "github.com/linzhengen/ddd-gin-admin/infrastructure/api/mock"
	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/domain/adapter"
	repo "github.com/linzhengen/ddd-gin-admin/domain/repository"
)

func BuildApiInjector() (*api.Injector, func(), error) {
	wire.Build(
		// mock.MockSet,
		api.InitGormDB,
		repo.RepoSet,
		api.InitAuth,
		api.InitCasbin,
		api.InitGinEngine,
		application.ServiceSet,
		handler.APISet,
		router.RouterSet,
		adapter.CasbinAdapterSet,
		api.InjectorSet,
	)
	return new(api.Injector), nil, nil
}
