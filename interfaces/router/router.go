package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/linzhengen/ddd-gin-admin/interfaces/handler"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
)

var _ IRouter = (*Router)(nil)

var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	Auth           auth.Author
	CasbinEnforcer *casbin.SyncedEnforcer
	DemoAPI        *handler.Demo
	LoginAPI       *handler.Login
	MenuAPI        *handler.Menu
	RoleAPI        *handler.Role
	UserAPI        *handler.User
	HealthAPI      *handler.HealthCheck
}

func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}
