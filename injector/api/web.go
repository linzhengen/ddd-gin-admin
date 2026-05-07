package api

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files/v2"
	"github.com/swaggo/swag"

	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/middleware"
	"github.com/linzhengen/ddd-gin-admin/app/interfaces/api/router"
	"github.com/linzhengen/ddd-gin-admin/configs"
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
		app.GET("/swagger/*any", func(c *gin.Context) {
			path := c.Param("any")
			// Gin includes the leading "/" in the catch-all param, strip it.
			if len(path) > 0 && path[0] == '/' {
				path = path[1:]
			}
			switch path {
			case "doc.json":
				doc, err := swag.ReadDoc()
				if err != nil {
					c.String(http.StatusInternalServerError, "failed to read swagger doc")
					return
				}
				c.String(http.StatusOK, doc)
			case "swagger-initializer.js":
				c.Header("Content-Type", "application/javascript; charset=utf-8")
				c.String(http.StatusOK, `window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "doc.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};`)
			default:
				http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerFiles.FS))).ServeHTTP(c.Writer, c.Request)
			}
		})
	}

	// Website
	if dir := configs.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	}

	return app
}
