package router

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/middleware"
)

// RegisterAPI register api group router
func (a *router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")
	g.Use(middleware.UserAuthMiddleware(a.auth,
		middleware.AllowPathPrefixSkipper("/api/v1/pub/login"),
	))

	g.Use(middleware.CasbinMiddleware(a.casbinEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/pub"),
	))

	g.Use(middleware.RateLimiterMiddleware())
	g.GET("health", a.healthHandler.Get)
	v1 := g.Group("/v1")
	{
		pub := v1.Group("/pub")
		{
			gLogin := pub.Group("login")
			{
				gLogin.GET("captchaid", a.loginHandler.GetCaptcha)
				gLogin.GET("captcha", a.loginHandler.ResCaptcha)
				gLogin.POST("", a.loginHandler.Login)
				gLogin.POST("exit", a.loginHandler.Logout)
			}

			gCurrent := pub.Group("current")
			{
				gCurrent.PUT("password", a.loginHandler.UpdatePassword)
				gCurrent.GET("user", a.loginHandler.GetUserInfo)
				gCurrent.GET("menutree", a.loginHandler.QueryUserMenuTree)
			}
			pub.POST("/refresh-token", a.loginHandler.RefreshToken)
		}

		gMenu := v1.Group("menus")
		{
			gMenu.GET("", a.menuHandler.Query)
			gMenu.GET(":id", a.menuHandler.Get)
			gMenu.POST("", a.menuHandler.Create)
			gMenu.PUT(":id", a.menuHandler.Update)
			gMenu.DELETE(":id", a.menuHandler.Delete)
			gMenu.PATCH(":id/enable", a.menuHandler.Enable)
			gMenu.PATCH(":id/disable", a.menuHandler.Disable)
		}
		v1.GET("/menus.tree", a.menuHandler.QueryTree)

		gRole := v1.Group("roles")
		{
			gRole.GET("", a.roleHandler.Query)
			gRole.GET(":id", a.roleHandler.Get)
			gRole.POST("", a.roleHandler.Create)
			gRole.PUT(":id", a.roleHandler.Update)
			gRole.DELETE(":id", a.roleHandler.Delete)
			gRole.PATCH(":id/enable", a.roleHandler.Enable)
			gRole.PATCH(":id/disable", a.roleHandler.Disable)
		}
		v1.GET("/roles.select", a.roleHandler.QuerySelect)

		gUser := v1.Group("users")
		{
			gUser.GET("", a.userHandler.Query)
			gUser.GET(":id", a.userHandler.Get)
			gUser.POST("", a.userHandler.Create)
			gUser.PUT(":id", a.userHandler.Update)
			gUser.DELETE(":id", a.userHandler.Delete)
			gUser.PATCH(":id/enable", a.userHandler.Enable)
			gUser.PATCH(":id/disable", a.userHandler.Disable)
		}
	}
}
