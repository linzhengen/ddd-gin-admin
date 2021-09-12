package api

import (
	"github.com/LyricTian/gzip"
	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/middleware"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/configs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitGinEngine(r router.Router) *gin.Engine {
	gin.SetMode(configs.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	prefixes := r.Prefixes()

	// Trace ID
	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Copy body
	app.Use(middleware.CopyBodyMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Access logger
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(prefixes...)))

	// Recover
	app.Use(middleware.RecoveryMiddleware())

	// CORS
	if configs.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	// GZIP
	if configs.C.GZIP.Enable {
		app.Use(gzip.Gzip(gzip.BestCompression,
			gzip.WithExcludedExtensions(configs.C.GZIP.ExcludedExtentions),
			gzip.WithExcludedPaths(configs.C.GZIP.ExcludedPaths),
		))
	}

	// Router register
	//nolint:errcheck
	r.Register(app)

	// Swagger
	if configs.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Website
	if dir := configs.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	}

	return app
}
