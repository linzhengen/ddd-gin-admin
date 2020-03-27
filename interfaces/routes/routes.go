package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/persistence/mysql"
	"github.com/linzhengen/ddd-gin-admin/interfaces/handler"
)

// RegisterRouter ...
func RegisterRouter(app *gin.Engine, mysqlRepo *mysql.Repositories) {
	g := app.Group("/api")
	v1 := g.Group("/v1")
	{
		// api/v1/users
		gUser := v1.Group("users")
		{
			uHandle := handler.NewUser(mysqlRepo.User)
			gUser.GET("", uHandle.Query)
			gUser.GET(":id", uHandle.Get)
			gUser.POST("", uHandle.Create)
			gUser.PUT(":id", uHandle.Update)
			gUser.DELETE(":id", uHandle.Delete)
			gUser.PATCH(":id/enable", uHandle.Enable)
			gUser.PATCH(":id/disable", uHandle.Disable)
		}
	}
}
