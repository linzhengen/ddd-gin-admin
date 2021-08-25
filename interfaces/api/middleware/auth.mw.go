package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/config"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/contextx"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/ginx"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
	"github.com/linzhengen/ddd-gin-admin/pkg/errors"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	ginx.SetUserID(c, userID)
	ctx := contextx.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

func UserAuthMiddleware(a auth.Author, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, config.C.Root.UserName)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), ginx.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config.C.IsDebugMode() {
					wrapUserAuthContext(c, config.C.Root.UserName)
					c.Next()
					return
				}
				ginx.ResError(c, errors.ErrInvalidToken)
				return
			}
			ginx.ResError(c, errors.WithStack(err))
			return
		}

		wrapUserAuthContext(c, userID)
		c.Next()
	}
}
