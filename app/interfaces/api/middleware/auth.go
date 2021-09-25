package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/contextx"
	"github.com/linzhengen/ddd-gin-admin/app/domain/valueobject/errors"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api"
	"github.com/linzhengen/ddd-gin-admin/configs"
	"github.com/linzhengen/ddd-gin-admin/pkg/auth"
	"github.com/linzhengen/ddd-gin-admin/pkg/logger"
)

func wrapUserAuthContext(c *gin.Context, userID string) {
	api.SetUserID(c, userID)
	ctx := contextx.NewUserID(c.Request.Context(), userID)
	ctx = logger.NewUserIDContext(ctx, userID)
	c.Request = c.Request.WithContext(ctx)
}

func UserAuthMiddleware(a auth.Author, skippers ...SkipperFunc) gin.HandlerFunc {
	if !configs.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapUserAuthContext(c, configs.C.Root.UserName)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, err := a.ParseUserID(c.Request.Context(), api.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if configs.C.IsDebugMode() {
					wrapUserAuthContext(c, configs.C.Root.UserName)
					c.Next()
					return
				}
				api.ResError(c, errors.ErrInvalidToken)
				return
			}
			api.ResError(c, errors.WithStack(err))
			return
		}

		wrapUserAuthContext(c, userID)
		c.Next()
	}
}