package main

import (
	"context"
	"os"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/persistence/mysql"

	"github.com/linzhengen/ddd-gin-admin/configs"

	"github.com/linzhengen/ddd-gin-admin/infrastructure/logger"
)

type options struct {
	SwaggerDir string
	Version    string
}

// Option option func.
type Option func(*options)

// SetSwaggerDir set swagger dir.
func SetSwaggerDir(s string) Option {
	return func(o *options) {
		o.SwaggerDir = s
	}
}

// SetVersion set app version.
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Init 应用初始化
func Init(ctx context.Context, opts ...Option) func() {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	err := configs.InitEnv()
	handleError(err)

	logger.Printf(ctx, "start server，run mode：%s，version：%s，进程号：%d", configs.Env().RunMode, o.Version, os.Getpid())

	loggerCall, err := InitLogger()
	handleError(err)
	db, storeCall, err := InitStore()
	handleError(err)
	mysqlRepo, err := mysql.NewRepositories(db)
	handleError(err)

	// 初始化HTTP服务
	httpCall := InitHTTPServer(ctx, mysqlRepo)
	return func() {
		if httpCall != nil {
			httpCall()
		}
		if loggerCall != nil {
			loggerCall()
		}
		if storeCall != nil {
			storeCall()
		}
	}
}
