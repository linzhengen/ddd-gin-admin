// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package injector

import (
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/interfaces/handler"

	// "github.com/linzhengen/ddd-gin-admin/infrastructure/api/mock"
	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/interfaces/router"

	"github.com/linzhengen/ddd-gin-admin/domain/adapter"
	repo "github.com/linzhengen/ddd-gin-admin/domain/repository"
)

// BuildInjector 生成注入器
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		// mock.MockSet,
		InitGormDB,
		repo.RepoSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		application.ServiceSet,
		handler.APISet,
		router.RouterSet,
		adapter.CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
