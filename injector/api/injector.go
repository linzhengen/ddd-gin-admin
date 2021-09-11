package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/application"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
)

func NewInjector(
	Engine *gin.Engine,
	Auth auth.Author,
	CasbinEnforcer *casbin.SyncedEnforcer,
	MenuBll application.Menu,
) *Injector {
	return &Injector{
		Engine:         Engine,
		Auth:           Auth,
		CasbinEnforcer: CasbinEnforcer,
		MenuBll:        MenuBll,
	}
}

type Injector struct {
	Engine         *gin.Engine
	Auth           auth.Author
	CasbinEnforcer *casbin.SyncedEnforcer
	MenuBll        application.Menu
}
