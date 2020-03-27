package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/persistence/mysql"

	"github.com/linzhengen/ddd-gin-admin/interfaces/routes"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/interfaces/middleware"

	"github.com/gin-gonic/gin"
	"github.com/linzhengen/ddd-gin-admin/infrastructure/logger"
)

// InitWeb init web.
func InitWeb(mysqlRepo *mysql.Repositories) *gin.Engine {
	cfg := configs.Env()
	gin.SetMode(cfg.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	apiPrefixes := []string{"/api/"}

	app.Use(middleware.TraceMiddleware(middleware.AllowPathPrefixNoSkipper(apiPrefixes...)))
	app.Use(middleware.LoggerMiddleware(middleware.AllowPathPrefixNoSkipper(apiPrefixes...)))
	app.Use(middleware.RecoveryMiddleware())

	routes.RegisterRouter(app, mysqlRepo)
	return app
}

// InitHTTPServer init http server
func InitHTTPServer(ctx context.Context, mysqlRepo *mysql.Repositories) func() {
	cfg := configs.Env()
	addr := fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      InitWeb(mysqlRepo),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Printf(ctx, "http server is starting，listen address：[%s]", addr)
		var err error
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf(ctx, err.Error())
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, cfg.HTTPShutdownTime)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, err.Error())
		}
	}
}
