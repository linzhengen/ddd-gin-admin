package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/domain/errors"
	"github.com/linzhengen/ddd-gin-admin/app/infrastructure/ginx"
	"github.com/linzhengen/ddd-gin-admin/configs"
)

func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := configs.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		if b, err := enforcer.Enforce(ginx.GetUserID(c), p, m); err != nil {
			ginx.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			ginx.ResError(c, errors.ErrNoPerm)
			return
		}
		c.Next()
	}
}
