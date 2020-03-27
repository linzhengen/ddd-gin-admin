package routes

import (
	"github.com/gin-gonic/gin"
)

// RegisterRouter ...
func RegisterRouter(app *gin.Engine) {
	g := app.Group("/api")
	v1 := g.Group("/v1")
	{
		// api/v1/users
		gUser := v1.Group("users")
		{
			gUser.GET("", h.Query)
			gUser.GET(":id", h.Get)
			gUser.POST("", cUser.Create)
			gUser.PUT(":id", cUser.Update)
			gUser.DELETE(":id", cUser.Delete)
			gUser.PATCH(":id/enable", cUser.Enable)
			gUser.PATCH(":id/disable", cUser.Disable)
		}
	}
}
